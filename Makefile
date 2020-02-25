ARTIFACT_DIRECTORY=/Users/leewaiho/Library/Application\ Support/Alfred/Alfred.alfredpreferences/workflows/user.workflow.FFE6060F-5DAF-44EC-91F2-4F8C26B22FBF
ARTIFACT_NAME=workflows

.PHONY: compile
compile:
	go build

.PHONY: install
install:
	@echo "编译程序中..." && go build -o $(ARTIFACT_DIRECTORY)/$(ARTIFACT_NAME) && echo "编译程序完成" && \
	 echo "安装配置文件中..." && cp build/workflows.json $(ARTIFACT_DIRECTORY)/workflows.json && echo "安装配置文件完成"

.PHONY: release
release:
	@bash scripts/release.sh