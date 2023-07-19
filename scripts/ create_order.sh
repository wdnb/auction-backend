#!/bin/bash

# 读取用户输入的订单信息
echo "请输入订单信息："
read order_info

# 调用 curl 命令提交订单信息
curl -X POST -H "Content-Type: application/json" -d "$order_info" http://localhost:8080/orders