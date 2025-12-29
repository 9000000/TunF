# TunF - Windows 11 Port Forwarding Utility

Ứng dụng proxy TCP hiện đại, nhẹ và mạnh mẽ dành cho Windows 11, được xây dựng bằng Wails v2 (Go & React).
Phục vụ nhu cầu thử nhiệm các dự án nội bộ cần kết nối ra internet nhanh chóng 

## ✨ Tính năng chính

- **Proxy TCP**: Chuyển hướng traffic từ cổng cục bộ sang bất kỳ địa chỉ đích nào.
- **Giao diện Win 11**: Thiết kế Glassmorphism (Mica) mượt mà, hỗ trợ giao diện tối.
- **Tự động mở Firewall**: Tự động quản lý luật tường lửa Windows (yêu cầu quyền Admin).
- **Lịch sử & Ghi nhớ**: Tự động lưu và gợi ý 10 cổng/địa chỉ hay dùng nhất.
- **Icon Khay hệ thống (System Tray)**:
  - Bật/Tắt proxy nhanh chóng.
  - Menu "Recent Connections" để chọn nhanh kết nối cũ.
  - Chạy ngầm khi đóng cửa sổ.

## 🚀 Hướng dẫn sử dụng

1. **Khởi chạy**: Mở ứng dụng `TunF`.
2. **Cấu hình**:
   - **Proxy Port**: Nhập cổng bạn muốn lắng nghe trên máy tính này (Ví dụ: `5678`).
   - **Target Address**: Nhập địa chỉ nguồn bạn muốn forward (Ví dụ: `localhost:1234`).
3. **Bật Proxy**: Nhấn nút **Start Proxy**.
   - Nếu bạn muốn thiết bị khác trong mạng cùng truy cập, hãy tích chọn **Auto-open Windows Firewall**.
4. **Khay hệ thống**: Khi ứng dụng đang chạy, bạn có thể click chuột phải vào icon ở khay hệ thống để:
   - Bật/Tắt proxy nhanh.
   - Chọn nhanh các kết nối trong quá khứ ở mục **Recent Connections**.
   - Hiện app hoặc Thoát hoàn toàn.

## 🛠️ Hướng dẫn Build ứng dụng

### 1. Yêu cầu môi trường

- **Go**: Phiên bản 1.18 trở lên.
- **Node.js & NPM**: Để build frontend.
- **Wails CLI**: Cài đặt bằng lệnh `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.

### 2. Chạy chế độ Development

Để chạy và chỉnh sửa code theo thời gian thực:

```bash
wails dev
```

### 3. Build file EXE chính thức

Để tạo file thực thi (`.exe`) gọn nhẹ và tối ưu:

```bash
wails build -clean -ldflags "-s -w"
```

File kết quả sẽ nằm trong thư mục `build/bin/`.

---

> [!IMPORTANT]
> Để sử dụng tính năng **Auto-open Firewall**, bạn cần chạy ứng dụng với quyền **Administrator**.

---

# 🇬🇧 English Instructions

## ✨ Key Features

- **TCP Proxy**: Forward traffic from a local port to any target address.
- **Windows 11 UI**: Sleek Glassmorphism (Mica) design, dark mode support.
- **Auto-open Firewall**: Automatically manages Windows Firewall rules (requires Admin privileges).
- **History & Favorites**: Automatically saves and suggests the 10 most used ports/addresses.
- **System Tray Icon**:
  - Quick Toggle Proxy ON/OFF.
  - "Recent Connections" menu for quick selection.
  - Runs in background when window is closed.

## 🚀 Usage Guide

1.  **Launch**: Open `TunF` application.
2.  **Configure**:
    - **Proxy Port**: Enter the local port to listen on (e.g., `5678`).
    - **Target Address**: Enter the target address to forward to (e.g., `localhost:1234`).
3.  **Start Proxy**: Click **Start Proxy**.
    - Check **Auto-open Windows Firewall** if you want other devices on the network to access it.
4.  **System Tray**: When running, right-click the tray icon to:
    - Quick Toggle Proxy.
    - Select from **Recent Connections**.
    - Show window or Quit.

## 🛠️ Build Instructions

### 1. Requirements

- **Go**: v1.18 or higher.
- **Node.js & NPM**: For frontend build.
- **Wails CLI**: Install via `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.

### 2. Development Mode

To run with live reload:

```bash
wails dev
```

### 3. Build Production EXE

To build a lightweight, optimized executable:

```bash
wails build -clean -ldflags "-s -w"
```

Output file will be in `build/bin/`.

---

> [!IMPORTANT]
> To use **Auto-open Firewall**, run the application as **Administrator**.
