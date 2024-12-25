<div align=center>
  <img width =256 src="./logo.jpg"/>
  <h1>AlphaLab</h1>
</div>

## Table of Contents
- [Introduction](#introduction)
- [Architecture Overview](#architecture-overview)
- [Features](#features)
	- [Key Features](#key-features)
	- [Small Features](#small-features)
- [Getting Started](#getting-started)
	- [Installation](#installation)
 	- [Development Roadmap](#development-roadmap)
- [Support and Contact](#support-and-contact)

---

## Introduction
**AlphaLab** is a comprehensive tool designed to manage laboratory operations effectively. It includes:
- A modern **frontend** built with **Vite** and **shadcn**, offering a clean and responsive user interface.
- A robust **backend** powered by **Golang** and **PocketBase** for seamless data management.
- A convenient **CLI** tool that simplifies software setup using **Docker**.

AlphaLab is open-source under the [Apache License](LICENSE).

---

## Architecture Overview
AlphaLab is divided into three main parts:

1. **Frontend (FE)**  
   - Built with [Vite](https://vitejs.dev/) for fast development.
   - Styled using [shadcn](https://shadcn.dev/) for a modern UI.

2. **Backend (BE)**  
   - Core logic and APIs built with **Golang**.
   - **PocketBase** serves as the lightweight database backend for managing lab data and user information.

3. **Command-Line Interface (CLI)**  
   - Streamlines setup and configuration tasks.
   - Requires **Docker** for containerized deployment.

---
## Features
### Key Features
1. **Lab Book Management**
   - Users can upload lab books in PDF format along with a description.
   - Files can be stored:
     - Locally
     - In **Google Cloud Storage** (requires user credentials)
     - In an **AWS S3 bucket** (requires user credentials)
   - Teachers are notified when a lab book is uploaded and can:
     - View the file and description.
     - Decide to **verify** or **decline** the submission.

2. **Schedule**
   - Manage lab schedules, experiment timelines, and team availability.
   - Calendar integration for event syncing.

3. **Resource Management**
   - Track equipment usage, inventory, and supplies.
   - Generate detailed reports for inventory audits.

### Small Features
1. **SMTP Mailer**
   - The Docker setup includes a small SMTP mailer for sending notifications.
   - Users must provide a domain for the mailer to work.

2. **Notification System**
   - Notifications are sent to teachers when a lab book is uploaded.
   - Configurable notification methods include:
     - Email (via the SMTP mailer)
     - LINE
     - Slack
     - Other methods (can be expanded based on integrations)

3. **Google Cloud Storage and AWS S3 Integration**
   - Users can choose where to store their uploaded files:
     - Google Cloud Storage (requires account credentials)
     - AWS S3 bucket (requires access credentials)
   - File storage preferences are configurable in system settings.

4. **User Management**
   - Role-based access controls (e.g., Admin, Researcher).
   - User activity tracking.

5. **System Settings**
   - Configure settings like time zones, storage, notification methods, and preferences.

---

## Getting Started

### Prerequisites
1. **Docker**  
   Ensure Docker is installed on your system. You can download it [here](https://www.docker.com/).

2. **Domain**  
   AlphaLab requires a domain to set up its SMTP server. You can purchase a domain from any provider, such as [Namecheap](https://www.namecheap.com/) or [GoDaddy](https://www.godaddy.com/).

---

## Installation

### Step 1: Clone the Repository
...


---
## Development Roadmap
### Frontend
- [ ] **Core UI Design and Functionality**
  - [ ] Login Page
  - [ ] Main page
  - [ ] Components
    - [ ] Dashboard
    - [ ] Schedule
    - [ ] Lab Book
      - [ ] Uplaod
      - [ ] Verify
      - [ ] History
    - [ ] Resources
    - [ ] Settings
    - [ ] Users
  - [ ] Implement role-based access features for different users (Admin, Teacher, User).


### Backend
- [ ] **Core API Development**
  - [ ] Implement APIs for user authentication, lab book uploads, and notifications.
  - [ ] Set up data models in PocketBase for users, lab books, schedules, and resources.

- [ ] **Integration and Scalability**
  - [ ] Integrate Google Cloud Storage and AWS S3 for file storage.
  - [ ] Add support for third-party notification systems (e.g., Slack, LINE).


### CLI
- [ ] **Setup and Configuration**
  - [ ] Develop scripts for deploying AlphaLab using Docker.
  - [ ] Simplify configuration options for storage and notifications.

- [ ] **Enhanced Features**
  - [ ] Add commands for managing users, resetting settings, and backing up data.
  - [ ] Enable live logs for monitoring backend and frontend performance.

- [ ] **Documentation and Automation**
  - [ ] Document all CLI commands with usage examples.
  - [ ] Automate deployment workflows using CI/CD pipelines.
"""

---

## Support and Contact
If you have any questions or issues:
- GitHub Issues: Report issues here
- Email: ...
