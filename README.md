<div align="center">
  
# social_network

</div>

<div align="center">

[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/Barsu5489/social_network)
</div>

This repository contains the codebase for a social network application. It includes both a backend (written in Go) and a frontend (built with Next.js), designed to provide users with a platform to connect, share posts, join groups, and participate in real-time conversations.

## ✨ Key Features

-   **User Authentication**: Secure user registration and login functionality.
-   **Profile Management**: Users can create and customize their profiles, controlling privacy settings.
-   **Post Creation**: Users can create and share posts with text and images, with privacy controls, comments, and likes.
-   **Follower System**:  Users can follow others and manage follow requests for private profiles.
-   **Group Management**: Discover and create groups, manage invitations, and group membership requests.
-   **Event Management**: Create, RSVP, and manage events within groups.
-   **Real-time Notifications**: Receive real-time notifications for follows, group invites, join requests, and new events.
-   **Real-time Chat**: Functionality for direct and group chats using WebSockets.

## 🧭 Table of Contents

-   [Installation](#-installation)
-   [Running the Project](#-running-the-project)
-   [Dependencies](#-dependencies)
-   [Contributing](#-contributing)
-   [License](#-license)
-   [Contact](#-contact)

## 🛠️ Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/Barsu5489/social_network.git
    cd social_network
    ```

2.  Navigate to the backend directory:

    ```bash
    cd backend
    ```

3.  Download the necessary Go modules:

    ```bash
    go mod tidy
    ```

4.  Navigate to the frontend directory:

    ```bash
    cd ../frontend
    ```

5.  Install the npm packages:

    ```bash
    npm install
    ```

## 🚀 Running the Project

1.  Start the backend server:

    ```bash
    cd backend
    go run server.go
    ```

2.  Start the frontend development server:

    ```bash
    cd frontend
    npm run dev
    ```

    Alternatively, to build and start the frontend in production mode:

    ```bash
    npm run build
    npm run start
    ```

## ⚙️ Dependencies

### Backend (Go)

-   `github.com/golang-migrate/migrate/v4`: Database migration tool.
-   `github.com/google/uuid`: Library for generating UUIDs.
-   `github.com/gorilla/mux`: HTTP request router and URL matcher.
-   `github.com/gorilla/sessions`: Provides cookie-based sessions.
-   `github.com/gorilla/websocket`: WebSocket implementation.
-   `github.com/mattn/go-sqlite3`: SQLite3 driver for Go.
-   `github.com/rs/cors`: Middleware for handling Cross-Origin Resource Sharing.
-   `golang.org/x/crypto`: Cryptographic primitives.

### Frontend (Next.js)

-   `next`: React framework for building web applications.
-   `react`: JavaScript library for building user interfaces.
-   `react-dom`: Serves as the entry point to the DOM or server-side rendering APIs of React.
-   `tailwindcss`: CSS framework.
-   `@radix-ui/react-*`: Set of primitive UI components.
-   `clsx`: Utility for constructing `className` strings conditionally and effectively.
-   `tailwind-merge`: Tool for resolving Tailwind CSS conflicts.
-   `lucide-react`: Collection of icons as React components.

## 🤝 Contributing

Contributions are welcome! Here's how you can contribute:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes and commit them.
4.  Submit a pull request.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📧 Contact

Maintainers:

-   [Barsu5489](https://github.com/Barsu5489)