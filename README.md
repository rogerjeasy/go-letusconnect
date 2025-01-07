
---

## 📂 `README.md`

# 🚀 Master's Program Networking & Collaboration Platform - Backend

This repository contains the backend implementation for the **Master's Program Networking & Collaboration Platform (letusconnect.com)**. It is a RESTful API built with **Golang** and the **Fiber** framework, providing endpoints for user authentication, networking, collaboration, and mentorship.

---

## 📋 Table of Contents

1. [Project Setup](#project-setup)
2. [Technologies Used](#technologies-used)
3. [Folder Structure](#folder-structure)
4. [Key Features](#key-features)
5. [Environment Variables](#environment-variables)
6. [Running the Server](#running-the-server)
7. [Database Setup](#database-setup)
8. [API Documentation](#api-documentation)
9. [Contributing](#contributing)
10. [License](#license)

---

## 🛠️ Project Setup

### 1. Clone the Repository

```bash
git clone https://github.com/rogerjeasy/go-letusconnect.git
cd go-letusconnect
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Create a `.env` File

Create a `.env` file in the root directory and add the following environment variables:

```env
PORT=8080

SMTP_HOST=your-smtp-host
SMTP_PORT=your-smtp-port
SENDER_EMAIL=your-email
SENDER_NAME=your-name
SENDER_PASS=your-email-password
FIREBASE_API_KEY=your-firebase-api-key
ENV_SERVICE_ACCOUNT_KEY=local
JSON_SERVICE_ACCOUNT_PATH=service-account.json
PUSHER_APP_ID=your-pusher-app-id
PUSHER_KEY=your-pusher-key
PUSHER_SECRET=your-pusher-secret
PUSHER_CLUSTER=your-pusher-cluster
```

### 4. Running the Server

```bash
go run main.go
```

---

## 💻 Technologies Used

- **Golang** - A statically typed, compiled programming language designed at Google.
- **Fiber** - An Express.js inspired web framework built on top of Fasthttp.
- **Firestore** - A flexible, scalable database for mobile, web, and server development from Firebase and Google Cloud.
- **Firebase Authentication** - A service that can authenticate users using only client-side code.
- **Pusher** - A real-time messaging service that allows you to add real-time features to your applications.

---

## 📁 Folder Structure

```bash

go-letusconnect
├── config                  # Configuration files
├── handlers                # Request handlers
├── middlewares             # Middleware functions
├── models                  # Data models
├── routes                  # API route definitions
├── services                # application logic and services
├── utils                   # Utility functions
├── .env                    # Environment variables
├── go.mod                  # Go module file
├── go.sum                  # Go sum file
├── main.go                 # Entry point of the application
└── README.md               # README file

```

---

## 🔑 Key Features

- **User Authentication** - Register, login, and logout users.
- **Profile Management** - Update user profile, change password, and reset password.
- **Networking Tools** - Connect with other users, view connections, and send connection requests.
- **Collaboration Tools** - Create, view, and join groups, post messages, and comment on posts.
- **Mentorship Platform** - Apply to be a mentor or mentee, view mentorship requests, and accept or reject requests.
- **Real-time messaging** - Receive real-time notifications for new connection requests, group invitations, and mentorship requests.
- **Email Notifications** - Send email notifications for new connection requests, group invitations, and mentorship requests.
- **Event Management** - Create, view, and join events, and view event details.
- **Job Board and Career Resources** - Post job listings, view job details, and apply for jobs.

---

## 🌐 Environment Variables

- **PORT** - The port number on which the server will run.
- **SMTP_HOST** - The SMTP host for sending email notifications.
- **SMTP_PORT** - The SMTP port for sending email notifications.
- **SENDER_EMAIL** - The email address of the sender for email notifications.
- **SENDER_NAME** - The name of the sender for email notifications.
- **SENDER_PASS** - The password of the sender's email address for email notifications.
- **FIREBASE_API_KEY** - The API key for Firebase Authentication.
- **ENV_SERVICE_ACCOUNT_KEY** - The environment variable for the service account key.
- **JSON_SERVICE_ACCOUNT_PATH** - The path to the service account key file.
- **PUSHER_APP_ID** - The Pusher app ID.
- **PUSHER_KEY** - The Pusher key.
- **PUSHER_SECRET** - The Pusher secret.
- **PUSHER_CLUSTER** - The Pusher cluster.

---

## 🚀 Running the Server

To run the server, execute the following command:

```bash
go run main.go
```

The server will start on the port specified in the `.env` file.

---

## 📦 Database Setup

The application uses Firestore as the database. You can set up Firestore by following the steps below:

1. Go to the [Firebase Console](https://console.firebase.google.com/).
2. Create a new project.
3. Go to the Firestore section and create a new Firestore database.
4. Add a new collection called `users`.
5. Add a new collection called `connections`.
6. Add a new collection called `groups`.

---

## 📖 API Documentation

API documentation is available via Postman collection or Swagger (if implemented).

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/rogerjeasy/go-letusconnect/issues).

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/awesome-feature`)
3. Commit your Changes (`git commit -m 'Add some awesome feature'`)
4. Open a Pull Request
---

## 📜 License

Distributed under the MIT License. See `LICENSE` for more information.

