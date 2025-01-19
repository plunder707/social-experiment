# Social-Experiment

**Social-Experiment** is a production-ready social media platform designed to facilitate user engagement and interactions through seamless content sharing. Built with a robust Go (Golang) backend and an intuitive Angular frontend, it allows users to register, log in, create, and share posts in real time, providing a dynamic and interactive social experience.

## Features

- **User Authentication**
  - Register new accounts
  - Secure login/logout functionality

- **Real-Time Posting**
  - Create and share posts instantly
  - Live updates without page refreshes

- **Responsive Design**
  - Optimized for desktops, tablets, and mobile devices

- **Interactive UI**
  - Clean and modern interface powered by Angular Material

## Technologies Used

- **Backend:**
  - Go (Golang)
  - Gin Web Framework
  - MongoDB
  - WebSockets for real-time communication

- **Frontend:**
  - Angular
  - Angular Material
  - RxJS

- **DevOps:**
  - Docker & Docker Compose for containerization

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.20 or higher)
- [Node.js](https://nodejs.org/en/download/) (version 18 or higher)
- [Angular CLI](https://angular.io/cli) (`npm install -g @angular/cli`)
- [Docker](https://www.docker.com/get-started) (optional, for containerization)
- [MongoDB](https://www.mongodb.com/try/download/community) (if not using Docker)

### Backend Setup

1. **Clone the Repository**

2. Edit .env and setup your server.

   ```bash
   git clone https://github.com/plunder707/social-experiment.git
   cd social-experiment
   go mod init social-experiment
   go mod tidy
   go run main.go

   
