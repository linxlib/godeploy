package deploy

import (
	"context"
	"fmt"
	"github.com/linxlib/godeploy/controllers/models"
	"github.com/saracen/fastzip"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"time"
)

var deployRoutine chan int = make(chan int)
var done chan bool = make(chan bool)

type Message struct {
	ID      int
	Message string
}

var messageRoutine chan Message = make(chan Message)

type Data struct {
	file *multipart.FileHeader
	svc  *models.Service
	db   *gorm.DB
}

var dataQueue map[int]*Data = make(map[int]*Data)
var running bool

func run() {
	running = true
	go func() {
		for {
			select {
			case id := <-deployRoutine:
				if data, ok := dataQueue[id]; ok {
					time.Sleep(time.Second * 1)
					messageRoutine <- Message{ID: id, Message: "开始部署"}
					f, err := data.file.Open()
					if err != nil {
						fmt.Println(data.file.Filename, data.file.Size)
						messageRoutine <- Message{ID: id, Message: err.Error()}
						continue
					}
					//time.Sleep(time.Second * 1)
					//TODO file hash check
					// TODO random folder
					messageRoutine <- Message{ID: id, Message: "解压到临时目录"}
					tmpEx := "./tmp"
					os.Mkdir(tmpEx, 0777)
					os.RemoveAll(tmpEx)
					reader, err := fastzip.NewExtractorFromReader(f, data.file.Size, tmpEx)
					if err != nil {
						messageRoutine <- Message{ID: id, Message: err.Error()}
						continue
					}
					//time.Sleep(time.Second * 1)
					for _, file := range reader.Files() {
						messageRoutine <- Message{ID: id, Message: file.Name}
					}
					err = reader.Extract(context.Background())
					if err != nil {
						messageRoutine <- Message{ID: id, Message: err.Error()}
						continue
					}
					err = f.Close()
					if err != nil {
						messageRoutine <- Message{ID: id, Message: err.Error()}
						continue
					}
					//time.Sleep(time.Second * 1)
					messageRoutine <- Message{ID: id, Message: "停止服务"}
					if !data.svc.Stop() {
						messageRoutine <- Message{ID: id, Message: "停止失败"}
						continue
					}
					//time.Sleep(time.Second * 1)
					messageRoutine <- Message{ID: id, Message: "覆盖已有服务文件"}
					if !data.svc.OverwriteFrom(tmpEx) {
						messageRoutine <- Message{ID: id, Message: "覆盖失败"}
						continue
					}
					//time.Sleep(time.Second * 1)
					messageRoutine <- Message{ID: id, Message: "重新启动服务"}
					if !data.svc.Start() {
						messageRoutine <- Message{ID: id, Message: "启动失败"}
						continue
					}
					//time.Sleep(time.Second * 1)
					data.svc.LastDeployTime = time.Now()
					if err := data.db.Save(data.svc).Error; err != nil {
						messageRoutine <- Message{ID: id, Message: "保存失败"}
						continue
					}
					messageRoutine <- Message{ID: id, Message: "部署成功"}
					done <- true
				} else {
					time.Sleep(time.Second * 1)
					messageRoutine <- Message{ID: id, Message: "未找到任务"}
					done <- true
				}
			}

		}

	}()
}

func CreateNewDeploy(file *multipart.FileHeader, svc *models.Service, db *gorm.DB) int {
	dataQueue[1] = &Data{
		file: file,
		svc:  svc,
		db:   db,
	}
	run()
	return 1
}

func StartDeploy(id int) (chan Message, chan bool) {
	if !running {
		run()
	}
	deployRoutine <- id
	return messageRoutine, done
}
