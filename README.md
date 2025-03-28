# SkyBox ☁️

[![Go Version](https://img.shields.io/badge/go-1.20%2B-00ADD8?logo=go)](https://golang.org/dl/)
[![License: GPL](https://img.shields.io/badge/license-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![AWS S3](https://img.shields.io/badge/storage-AWS_S3-FF9900?logo=amazon-aws)](https://aws.amazon.com/s3/)
[![MongoDB Atlas](https://img.shields.io/badge/database-MongoDB_Atlas-47A248?logo=mongodb)](https://www.mongodb.com/atlas/database)

<!---toc:start-->
- [SkyBox ☁️](#skybox-️)
  - [Features ✨](#features)
    - [Core Features](#core-features)
    - [Advanced Features](#advanced-features)
  - [Technology Stack 🛠️](#technology-stack-🛠️)
    - [Backend](#backend)
    - [Frontend](#frontend)
  - [Getting Started 🚀](#getting-started-🚀)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
      - [Install from Source](#install-from-source)
      - [Deployment](#deployment)
  - [Contributors](#contributors)
<!--toc:end-->

SkyBox is a secure, scalable cloud storage solution inspired by Google Drive and Dropbox, built with Go, AWS, and MongoDB. It offers file storage, sharing, and synchronization across devices.

> [!NOTE]
> This is the final group project from two courses in VNUHCM - University of Science - CS422 - Software Analysis and Design.

## Features ✨

### Core Features

- **File Management**: Upload, download, organize files and folders
- **User Authentication**: Secure signup/login with JWT
- **File Sharing**: Share files/folders with other users
- **Versioning**: Keep track of file versions

### Advanced Features

- **Chunked Uploads**: Support for large files
- **Real-time Sync**: WebSocket-based file synchronization
- **Search**: Full-text search across your files
- **Trash System**: Recover deleted files within retention period

## Technology Stack 🛠️

### Backend

- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: MongoDB
- **Object Storage**: AWS S3

### Frontend

- React

## Getting Started 🚀

### Prerequisites

- Go 1.20+
- MongoDB 6.0+
- AWS Account with S3 access
- Node.js (for frontend)

### Installation

#### Install from Source

1. Clone the repository:

```bash
git clone https://github.com/skybox-apcs/skybox-backend.git
cd skybox-backend
```

2. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies:

```bash
go mod download
```

4. Run the application:

```bash
go run .\cmd\server\main.go
```

5. Access the application:

```bash
http://localhost:8080
```

#### Deployment

TBA

## Contributors

The project could not have been completed without these developers!

- 22125050 - Nguyễn Thanh Phước Lộc
  - <ntploc22@apcs.fitus.edu.vn>
- 22125068 - Trương Chí Nhân
  - <tcnhan22@apcs.fitus.edu.vn>
- 22125076 - Nguyễn Hoàng Phúc
  - <nhphuc221@apcs.fitus.edu.vn>
- 22125115 - Ngô Hoàng Tuấn
  - <nhtuan22@apcs.fitus.edu.vn>

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](https://github.com/skybox-apcs/skybox-backend/blob/KAN-1/LICENSE) file for details.
