# P2K16 Development Environment Setup Guide

This comprehensive guide provides step-by-step instructions for setting up and running the P2K16 project. It covers prerequisites, installation procedures, environment configuration, and application deployment for both Linux (Ubuntu/Debian) and macOS platforms. Each tool and process is thoroughly documented to ensure a smooth development experience.

---

Quick start: see docs/LOCAL_DEV.md for an end-to-end local setup (macOS + fish) including Docker/Postgres, Flyway, and running both apps.

---

# ⚠️ Python Version Compatibility Warning

> **Important:**
> For best compatibility with scientific Python packages and to avoid installation errors (such as `pkgutil.ImpImporter` missing in Python 3.12+), it is strongly recommended to use **Python 3.11** for development. 
> 
> - If you are on macOS, install with:
>   ```sh
>   brew install python@3.11
>   ```
> - Create your virtual environment with:
>   ```sh
>   python3.11 -m venv env
>   source env/bin/activate.fish  # or activate for your shell
>   ```
> - Then proceed with `pip install --upgrade pip setuptools wheel` and `pip install -r requirements.txt`.
> 
> Using Python 3.12+ may result in errors due to legacy dependencies. If you encounter such errors, switch to Python 3.11.

---

## 🛠️ Prerequisites

The following tools are required for P2K16 development:

- **Python virtualenv**: Creates isolated Python environments for dependency management. Install using `pip3 install virtualenv` or your system's package manager.
- **Docker**: Containerization platform for running services. Refer to the [official Docker installation guide](https://docs.docker.com/get-docker/).
- **Node Version Manager (nvm)**: Manages multiple Node.js versions. Installation instructions available at the [nvm repository](https://github.com/nvm-sh/nvm).
- **PostgreSQL client & libraries**: Provides database connectivity and development headers. Includes the `psql` command-line interface and libraries required for Python/Java integration.
- **Java Runtime Environment (JRE)**: Required for executing Flyway database migrations.

---

## 📦 Installation

### Ubuntu/Debian

Execute the following commands to install all required dependencies:

```sh
sudo apt update
sudo apt install python3-pip python3-venv docker.io nvm postgresql-client-common postgresql-client libpq-dev default-jre
pip3 install virtualenv
```

### macOS

Install dependencies using [Homebrew](https://brew.sh/):

```sh
brew install python@3 virtualenv docker nvm libpq postgresql openjdk
```

#### 🐧 Docker Backend on macOS: Colima

On macOS, you need a Docker backend to run containers. You can use either Docker Desktop (official, GUI-based) or Colima (lightweight, open source, CLI-based).

**To use Colima (recommended for developers):**
```sh
brew install colima
colima start
```
Then install the Docker CLI if you haven’t already:
```sh
brew install docker
```
Test Docker is working:
```sh
docker info
```
You should see Docker engine info. If you get an error, make sure Colima is running.


#### Configure PostgreSQL client tools

Homebrew does not automatically link PostgreSQL client tools to the system PATH. Add them manually:

For fish shell:
```fish
echo 'set -gx PATH /opt/homebrew/opt/libpq/bin $PATH' >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish
```
For bash/zsh:
```sh
echo 'export PATH="/opt/homebrew/opt/libpq/bin:$PATH"' >> ~/.zprofile
source ~/.zprofile
```

#### Configure JAVA_HOME environment variable

If required, set the JAVA_HOME environment variable:

For fish shell:
```fish
echo 'set -gx JAVA_HOME (brew --prefix openjdk)/libexec/openjdk.jdk/Contents/Home' >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish
```
For bash/zsh:
```sh
echo 'export JAVA_HOME="$(brew --prefix openjdk)/libexec/openjdk.jdk/Contents/Home"' >> ~/.zprofile
source ~/.zprofile
```

---

## ⚙️ Environment Setup

The project provides custom tools in the `bin/` directory. Before executing any project commands, source the appropriate setup script for your shell environment:

**For bash/zsh users:**
```sh
. .settings.sh
```

**For fish shell users:**
```fish
source .settings.fish
```

This configuration establishes the correct PATH and environment variables required for the project.

---

## 🐳 Docker Compose

This project utilizes Docker Compose for database service management. Docker Compose is a tool for defining and orchestrating multi-container Docker applications.

The available command varies based on your Docker installation: `docker-compose` (legacy standalone tool) or `docker compose` (modern Docker CLI plugin).

### Installation Instructions

**Ubuntu/Debian:**
- Docker installations following official documentation typically include the `docker compose` plugin.
- If unavailable, install the plugin manually:
  ```sh
  sudo apt install docker-compose-plugin
  ```

**macOS:**
- Docker Desktop includes Docker Compose as `docker compose` by default.
- Homebrew users can install the standalone version:
  ```sh
  brew install docker-compose
  ```

Verify installation by checking the available command:
```sh
docker compose version || docker-compose version
```

---

## 🗄️ Database Setup

1. **Start the database service using Docker Compose:**
   ```sh
   cd docker/p2k16
   docker-compose up -d
   cd -
   ```

2. **Initialize the database schema:**
   ```sh
   psql -U postgres -f database-setup.sql
   ```

**Database Credentials:** The default database user is `postgres` with password `postgres`.

---

## 🌐 Running the Application

Launch the web application using the provided script:

```sh
p2k16-run-web
```

**Troubleshooting:** If the application fails to start, verify that all required dependencies are properly installed. You may need to manually adjust `requirements.txt` if version conflicts occur. The `p2k16-update-requirements` script may resolve some dependency issues.

**Access:** Upon successful startup, the application will be accessible at [http://localhost:5000/](http://localhost:5000/).

---

## 👤 Default User Accounts

The application includes the following pre-configured user accounts for testing:

- **Administrator:** Username: `super`, Password: `super`
- **Standard User:** Username: `foo`, Password: `foo` (limited privileges)

---

## 💳 Account Payment Configuration

To configure an account as a paying member, manually insert a payment record into the database. Replace `2` with the target user ID:

```sql
INSERT INTO stripe_payment(created_at, created_by, updated_at, updated_by, stripe_id, start_date, end_date, amount, payment_date)
VALUES(now(), 2, now(), 2, 'fake-stripe-id-2', now(), now() + interval '10 years', 500.00, now());
```

---

## 🛠️ Development Information

## 📦 Project Dependencies

The P2K16 application is built using the following core technologies:

- **Flask**: Python web framework providing the graphical user interface ([Documentation](http://flask.pocoo.org))
- **PostgreSQL**: Relational database management system serving as the primary data store
- **SQLAlchemy**: Python Object-Relational Mapping (ORM) library for database interactions

## ⚡ Auto-generated Files

JavaScript files located in `/src/p2k16/web/static` are automatically generated by corresponding Python modules. For example, `door-data-service.js` is generated by `door_blueprint.py`.

## 🗃️ Database Schema Management

The project uses [Flyway](https://flywaydb.org) for database schema version control and migrations. After setting up the database, execute migrations with:

```sh
flyway migrate
```

**Creating New Migrations:** To modify the database schema, create a new file in the `migrations/` directory following the naming convention: `V001.NNN__descriptive_comment.sql`. Note that concurrent schema changes by multiple developers will result in merge conflicts.

---

## 📝 Development Roadmap

The following items represent planned enhancements and known issues:

- ✨ **UI Enhancement**: Implement word completion functionality for the Add Badge text field
- 🧹 **User Experience**: Clear text field automatically after badge addition on user profile
- ✅ **Feedback**: Display success confirmation message after badge addition on user profile
- 🔍 **Data Integrity**: Resolve duplicate badge name validation issues
- ☁️ **Visualization**: Implement word cloud or similar visualization for badges on Bitraf front page to enhance engagement
- 🗑️ **User Control**: Add delete functionality for self-created badges
- 📏 **Layout**: Implement proper line length constraints for badges on user profile
- 🔢 **Database Optimization**: Replace BIGSERIAL with BIGINT on version tables
- 🔒 **Data Protection**: Implement field update restrictions for sensitive fields like Account.username (consider SQLAlchemy event systems: http://docs.sqlalchemy.org/en/latest/orm/events.html)
- 🔄 **System Reliability**: Implement state persistence for tools during system reboots using retained MQTT messages

---

# 🏅 Badge System Documentation

## 🎯 System Overview

The badge system is a comprehensive recognition and authorization framework designed to serve multiple organizational objectives:

- 🛡️ **Safety Compliance**: Enforce mandatory training requirements for operating potentially dangerous equipment
- 📚 **Educational Incentives**: Provide recognition and motivation for course instructors and participants
- 🔍 **Skill Discovery**: Enable community members to identify expertise and available knowledge within the organization
- 🎉 **Community Engagement**: Encourage active participation and ongoing involvement in organizational activities

## 📖 Operational Framework

Badges serve as digital credentials that convey information about a user's skills, contributions, and authorizations. While badges carry no monetary value, they provide significant social recognition and functional access control.

The system implements a multi-tiered authorization model:
- **Peer-awarded badges**: Recognition that can only be granted by other community members (karma badges)
- **Authority-restricted badges**: Credentials that require specific authorization roles (e.g., course instructor privileges)
- **Cumulative badges**: Recognition that can be awarded multiple times to reflect ongoing contributions

## 🏷️ Badge Classification System

### 🔧 Equipment Competency Badges
*Required for operating specific tools that mandate safety training*
- Laser Cutter Operator
- CNC Machine Operator  
- Lathe Operation Certified

### 🛠️ Technical Skill Badges
*Recognition of expertise in various technical domains*
- Laser Cutting Proficiency
- Woodworking Expertise
- Metalworking Skills
- Soldering Techniques
- PCB Design (KiCAD, Eagle)
- Surface Mount Technology (SMT)
- PCB Etching Processes
- Oscilloscope Operation

### ⭐ Community Contribution Badges
*Recognition for ongoing organizational support*
- Facility Maintenance Volunteer
- Workshop Event Organizer (Dugnader)
- Infrastructure Project Contributor

### 💡 Special Recognition Categories
- **Innovation/Initiative Badge**: Awarded for significant contributions to organizational improvement (equipment repair, educational programs, system development)
- **Professional Certification**: Recognition of relevant professional qualifications (Professional Programmer, Licensed Electrician, Certified Carpenter)

---

## 🛠️ Troubleshooting Python Package Installation

If you encounter errors like `AttributeError: cython_sources`, `pkgutil.ImpImporter`, or errors mentioning `canonicalize_version()` or `strip_trailing_zero` when running `pip install -r requirements.txt` or `pip install -e web`, follow these steps:

1. **Upgrade pip, setuptools, wheel, cython, and packaging in your virtual environment:**
   ```sh
   pip install --upgrade pip setuptools wheel cython packaging
   ```
2. **Then install project requirements:**
   ```sh
   pip install -r requirements.txt
   ```

> **Note:** The `requirements.txt` now pins modern, compatible versions of `setuptools` and `packaging` to avoid these errors. If you still encounter issues, ensure your environment is using the correct Python version and all build tools are up to date.

This resolves most build and compatibility issues with legacy scientific Python packages on modern Python versions.

If you still encounter errors, check that you are using Python 3.11 as recommended above, and consider updating any problematic packages in `requirements.txt` to their latest versions.

---

## ☕ Troubleshooting Flyway and Java Runtime

If you see an error like:

> The operation couldn’t be completed. Unable to locate a Java Runtime.

This means Flyway (or other Java-based tools) cannot find a Java Runtime Environment (JRE).

### Solution (macOS):
1. Install OpenJDK via Homebrew:
   ```sh
   brew install openjdk
   ```
2. Add OpenJDK to your PATH and set JAVA_HOME (for fish shell):
   ```fish
   echo 'set -gx JAVA_HOME (brew --prefix openjdk)/libexec/openjdk.jdk/Contents/Home' >> ~/.config/fish/config.fish
   echo 'set -gx PATH $JAVA_HOME/bin $PATH' >> ~/.config/fish/config.fish
   source ~/.config/fish/config.fish
   ```
   For bash/zsh:
   ```sh
   echo 'export JAVA_HOME="$(brew --prefix openjdk)/libexec/openjdk.jdk/Contents/Home"' >> ~/.zprofile
   echo 'export PATH="$JAVA_HOME/bin:$PATH"' >> ~/.zprofile
   source ~/.zprofile
   ```

### Solution (Linux):
1. Install OpenJDK:
   ```sh
   sudo apt install default-jre
   ```
2. If needed, set JAVA_HOME:
   ```sh
   export JAVA_HOME="/usr/lib/jvm/default-java"
   export PATH="$JAVA_HOME/bin:$PATH"
   ```

After installing and configuring Java, re-run your Flyway or database migration command.

---

## 🗄️ Flyway Database Connection Configuration

If you see errors like:

> ERROR: Unable to connect to the database. Configure the url, user and password!

This means Flyway cannot find the database connection details. You must provide them via a configuration file or command-line arguments.

### Example command-line usage:
```sh
flyway -url=jdbc:postgresql://localhost:5432/postgres -user=postgres -password=postgres migrate
```

### Using a configuration file:
Create a file named `flyway.conf` (or use an existing one, e.g., in `infrastructure/` or `docker/p2k16/`) with contents like:
```
flyway.url=jdbc:postgresql://localhost:5432/postgres
flyway.user=postgres
flyway.password=postgres
```

Then run:
```sh
flyway migrate
```

> **Note:**
> - Ensure your database is running and accessible at the specified host and port.
> - The default credentials are `postgres`/`postgres` unless you have changed them.
> - You may need to adjust the host, port, or database name to match your environment.
