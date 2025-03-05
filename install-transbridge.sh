#!/bin/bash

# 获取脚本所在目录的绝对路径
SCRIPT_DIR=$(cd $(dirname $0) && pwd)

# 定义变量
SERVICE_NAME="transbridge-linux-amd64"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"
WORKING_DIRECTORY=$SCRIPT_DIR
EXEC_START="$WORKING_DIRECTORY/$SERVICE_NAME"
LOG_FILE="$WORKING_DIRECTORY/$SERVICE_NAME.log"
CONFIG_FILE="$WORKING_DIRECTORY/config_transbrige.yml"

# 创建 systemd 服务单元文件内容
SERVICE_CONTENT="[Unit]
Description=TransBridge Translation Service
After=network.target

[Service]
Type=simple
WorkingDirectory=$WORKING_DIRECTORY
ExecStart=$EXEC_START -config $CONFIG_FILE
Restart=on-failure
StandardOutput=append:$LOG_FILE
StandardError=append:$LOG_FILE

[Install]
WantedBy=multi-user.target"

# 检查必要文件是否存在
if [ ! -x "$EXEC_START" ]; then
    echo "错误: 可执行文件 $EXEC_START 不存在或不可执行"
    exit 1
fi

if [ ! -f "$CONFIG_FILE" ]; then
    echo "错误: 配置文件 $CONFIG_FILE 不存在"
    exit 1
fi

# 创建日志文件（如果不存在）
touch "$LOG_FILE"
chmod 644 "$LOG_FILE"

# 创建服务单元文件
echo "创建服务单元文件 $SERVICE_FILE"
echo "$SERVICE_CONTENT" | sudo tee $SERVICE_FILE > /dev/null

# 重新加载 systemd 配置
echo "重新加载 systemd 配置"
sudo systemctl daemon-reload

# 启动并启用服务
echo "启动 $SERVICE_NAME 服务"
sudo systemctl start $SERVICE_NAME

echo "启用 $SERVICE_NAME 服务在启动时自动运行"
sudo systemctl enable $SERVICE_NAME

echo "--------------------------------"
echo "$SERVICE_NAME 服务已安装并启动"
echo "日志文件位置: $LOG_FILE"
echo "可以使用以下命令管理服务："
echo "  启动: sudo systemctl start $SERVICE_NAME"
echo "  停止: sudo systemctl stop $SERVICE_NAME"
echo "  重启: sudo systemctl restart $SERVICE_NAME"
echo "  状态: sudo systemctl status $SERVICE_NAME"
echo "  查看日志: tail -f $LOG_FILE"