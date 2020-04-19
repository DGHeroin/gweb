package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/kardianos/service"
    logger "github.com/sirupsen/logrus"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"
    "runtime"
)

var (
    path = flag.String("d", "./", "path")
    port = flag.Int("p", 8080, "port")
)

type AppConfig struct {
    Dir string
    Port int
}

type program struct {}

func (p *program) Start(s service.Service) error {
    go p.run()
    return nil
}

func (p *program) run() {
    startServe()
}

func (p *program) Stop(s service.Service) error {
    return nil
}

func main() {
    svcConfig := &service.Config{
        Name:        "gweb",                              //服务显示名称
        DisplayName: "gweb Service",                      //服务名称
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
        action := os.Args[1]
        switch action {
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
        case "config":
            setConfig()
            return
        }
        return
    }

    err = s.Run()
    if err != nil {
        logger.Error(err)
    }
}
func setConfig() {
    os.Args = os.Args[1:]
    flag.Parse()
    cfg := AppConfig{}
    cfg.Dir = *path
    cfg.Port = *port
    data, _ := json.Marshal(&cfg)
    cfgPath := fmt.Sprintf("%v/.gweb.json", getPath())
    if err := ioutil.WriteFile(cfgPath, data, 0600); err != nil {
        logger.Error(err)
    }
    logger.Println("写配置成功:", getPath(), string(data))
}
func getConfig() (*AppConfig, error) {
    cfg := &AppConfig{
        Dir:  "./",
        Port: 8080,
    }
    cfgPath := fmt.Sprintf("%v/.gweb.json", getPath())
    data, err := ioutil.ReadFile(cfgPath)
    if err != nil {
        cfg := AppConfig{
            Dir:  "./",
            Port: 8080,
        }
        data, err := json.Marshal(&cfg)
        ioutil.WriteFile(cfgPath, data, 0600)
        return nil, err
    }
    err = json.Unmarshal(data, &cfg)
    return cfg, err
}

func getPath() string {
    switch runtime.GOOS {
    case "linux":
        os.Mkdir("/etc/gweb/", 0600)
        return "/etc/gweb/"
    default:
        dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            panic(err)
        }
        return dir
    }
}

func startServe() {
    cfg, err := getConfig()
    if err != nil {
        logger.Error(err)
        return
    }
    http.Handle("/", http.FileServer(http.Dir(cfg.Dir)))
    err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
    if err != nil {
        logger.Error(err)
    }
}
