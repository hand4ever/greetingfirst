.PHONY: rundev
rundev:
	@echo "本地运行 dev"
	go run main.go


.PHONY: buildqa
buildqa:
	@echo "测试服qa，部署到 111.333.222.444 服务器上"
	@echo "开始编译。。。"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \
	 -o bin/api main.go
	@echo "编译完成！"
	@echo "传输中。。。"
	@scp -O bin/api root@111.333.222.444:/opt/src/main
	@echo "传输完成！"
	@echo "重启服务中。。。"
	@ssh root@111.333.222.444  "supervisorctl restart xxxx"
	@echo "完成！"


