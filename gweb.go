package main
import (
    "flag"
    "fmt"
    "github.com/kardianos/service"
    "net/http"
    "os"
)

var (
    path = flag.String("d", "./", "path")
    port = flag.Int("p", 8080, "port")
)

var logger = service.ConsoleLogger

type program struct {
}

func (p *program) Start(s service.Service) error {
    go p.run()
    return nil
}

func (p *program) run() {
    loginFn()
}

func (p *program) Stop(s service.Service) error {
    return nil
}

func main() {
    svcConfig := &service.Config{
        Name:        "gweb",                 //服务显示名称
        DisplayName: "gweb Service",         //服务名称
        Description: "Anther other static file service.", //服务描述
    }

    prg := &program{}
    s, err := service.New(prg, svcConfig)
    if err != nil {
        logger.Error(err)
    }

    if err != nil {
        logger.Error(err)
    }

    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "install":
            s.Install()
            logger.Info("服务安装成功!")
            s.Start()
            logger.Info("服务启动成功!")
            break
        case "start":
            s.Start()
            logger.Info("服务启动成功!")
            break
        case "stop":
            s.Stop()
            logger.Info("服务关闭成功!")
            break
        case "restart":
            s.Stop()
            logger.Info("服务关闭成功!")
            s.Start()
            logger.Info("服务启动成功!")
            break
        case "remove":
            s.Stop()
            logger.Info("服务关闭成功!")
            s.Uninstall()
            logger.Info("服务卸载成功!")
            break
        case "status":
            if st, err := s.Status(); err != nil {
                logger.Error(err)
                return
            } else {
                logger.Infof("%s %v", s, st)
            }
        case "dummy":
            return
        }

        return
    }
    err = s.Run()
    if err != nil {
        logger.Error(err)
    }
}

func loginFn() {
    flag.Parse()
    http.Handle("/", http.FileServer(http.Dir(*path)))

    err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
    if err != nil {
        fmt.Println(err)
    }
}