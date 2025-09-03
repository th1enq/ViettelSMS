#!/bin/bash

# Đọc thông tin từ file .env nếu tồn tại
if [ -f ".env" ]; then
    # Đọc file .env và export các biến, xử lý dấu ngoặc kép
    while IFS='=' read -r key value; do
        # Bỏ qua dòng trống và comment
        [[ -z "$key" || "$key" =~ ^[[:space:]]*# ]] && continue
        # Loại bỏ khoảng trắng thừa
        key=$(echo "$key" | xargs)
        value=$(echo "$value" | xargs)
        # Export biến
        export "$key"="$value"
    done < .env
else
    echo "File .env không tồn tại. Vui lòng chạy register.sh trước."
    exit 1
fi

# Kiểm tra các biến cần thiết
if [ -z "$SERVER_ID" ] || [ -z "$INTERVAL_TIME" ]; then
    echo "Thiếu thông tin SERVER_ID hoặc INTERVAL_TIME trong file .env"
    exit 1
fi

# Hàm gửi metrics
send_metrics() {
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%S.000Z")
    
    echo "Gửi metrics"

    RESPONSE=$(curl -s -X POST http://$HOST_IP:80/api/v1/monitoring \
        -H "Content-Type: application/json" \
        -d "{
            \"server_id\": \"$SERVER_ID\",
            \"interval_time\": $INTERVAL_TIME,
            \"timestamp\": \"$timestamp\"
        }")

    echo "Response: $RESPONSE"
    
    
    return 0
}

# Hàm chính
main() {
    echo "Bắt đầu gửi metrics cho server $SERVER_ID với interval ${INTERVAL_TIME}s"
    
    # Kiểm tra jq có sẵn không
    if ! command -v jq &> /dev/null; then
        echo "jq chưa được cài đặt. Vui lòng cài đặt jq để parse JSON."
        exit 1
    fi
    
    # Gửi metrics đầu tiên ngay lập tức
    send_metrics
    
    
    # Lặp vô hạn với interval time
    while true; do
        sleep $INTERVAL_TIME
        send_metrics
        
        # Nếu gửi thất bại nhiều lần, có thể cần dừng
        if [ $? -ne 0 ]; then
            echo "Gặp lỗi khi gửi metrics. Thử lại sau ${INTERVAL_TIME}s..."
        fi
    done
}

# Xử lý tín hiệu để dừng script một cách graceful
trap 'echo "Dừng gửi metrics..."; exit 0' SIGINT SIGTERM

# Chạy hàm chính
main "$@"