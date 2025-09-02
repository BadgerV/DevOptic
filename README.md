# DevOptic
### Enterprise DevOps Monitoring & Management Platform

**Centralized visibility and control for microservices infrastructure at scale.**

---

## 🎯 **Problem Statement**

Managing 70+ microservices across multiple staging servers creates operational complexity. Development and DevOps teams waste valuable time jumping between GitLab CI/CD, Docker dashboards, service logs, and monitoring tools to get a complete picture of system health.

**DevOptic eliminates this cognitive overhead** by providing unified visibility, automated health checks, and streamlined operations management in a single, intuitive interface.

---

## ⚡ **Key Features**

### **🔐 Role-Based Access Control (RBAC)**
- **Super Admin** creates custom permissions and roles
- **Granular Permission System** - assign specific permissions to roles
- **User Management** - secure role assignment and access control
- **Audit Trail** - track all administrative actions

### **🔗 GitLab Integration**
- **Service Registration** - register repositories with name and ID mapping
- **Service Classification** - distinguish between microservices and macroservices
- **Dependency Management** - handle services that depend on other builds
- **Pipeline Visibility** - monitor CI/CD status across all registered services

### **📊 Health Monitoring**
- **Automated Endpoint Monitoring** - periodic health checks for all registered services
- **Intelligent Alerting** - email notifications when services go down
- **Real-time Status Dashboard** - instant visibility into service health
- **Scalable Architecture** - designed to handle enterprise-level microservices

---

## 🛠️ **Technology Stack**

| Component | Technology | Version |
|-----------|------------|---------|
| **Frontend** | React.js | 19.1.1 |
| **Backend** | Golang | 1.24.5 |
| **Database** | PostgreSQL | Latest |
| **Containerization** | Docker | Latest |
| **Orchestration** | Docker Compose | - |

### **Architecture Highlights**
- **Microservices-Ready Design** - built with future microservices migration in mind
- **Cloud-Native Patterns** - containerized for consistency across environments
- **High-Performance Backend** - Golang ensures efficient concurrent monitoring
- **Responsive UI** - React.js provides real-time dashboard updates

---

## 🚀 **Quick Start**

### **Prerequisites**
- Docker and Docker Compose installed
- Network access for service monitoring
- PostgreSQL connection available

### **Installation**
```bash
# Clone the repository
git clone https://github.com/[your-username]/devoptic.git

cd devoptic

# Start all services
docker-compose up -d

# Access the platform
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### **Network Access**
The platform can be accessed across your network via IP address for team collaboration.

---

## 📈 **Current Status**

**✅ MVP Features Complete:**
- RBAC system with permissions, roles, and user management
- GitLab repository registration and service classification
- Automated endpoint health monitoring
- Email alerting system for service failures
- Docker containerization with compose setup

**🔄 In Development:**
- Enhanced metrics collection and visualization
- Kubernetes integration for orchestration management
- Microservices architecture migration
- Advanced alerting and notification systems

---

## 🎯 **Use Cases**

### **For DevOps Engineers**
- Monitor 70+ microservices from a single dashboard
- Quickly identify failing services and their dependencies
- Streamline incident response with centralized visibility
- Manage team access with granular RBAC controls

### **For Development Teams**
- Track deployment status across multiple projects
- Understand service dependencies and build relationships
- Reduce time spent switching between monitoring tools
- Get automated notifications for service issues

### **For Engineering Managers**
- Gain enterprise-wide visibility into service health
- Track team productivity and deployment frequency
- Ensure proper access controls and security compliance
- Monitor overall infrastructure reliability metrics

---

## 🗺️ **Roadmap**

### **Phase 1: Enhanced Monitoring** (Q1 2025)
- Advanced metrics collection (CPU, memory, response times)
- Custom alerting rules and notification channels
- Service dependency mapping and visualization

### **Phase 2: Microservices Migration** (Q2 2025)
- Decompose monolithic backend into focused microservices
- Implement service mesh patterns for inter-service communication
- Enhanced scalability and deployment flexibility

### **Phase 3: Enterprise Integration** (Q3-Q4 2025)
- Kubernetes cluster management and monitoring
- Multi-cloud support (AWS, Azure integration)
- Advanced analytics and capacity planning
- Integration with HashiCorp Vault for secrets management

---

## 🤝 **Contributing**

Currently in private development with planned open-source release. Interested in enterprise DevOps monitoring solutions? Connect with me on [LinkedIn](https://www.linkedin.com/in/faozan-segunmaru-258502200).

---

## 📧 **Contact**

**Faozan Segunmaru**  
DevOps Engineer | Enterprise Infrastructure Specialist  
Building the future of microservices monitoring

---


*Built with ❤️ for DevOps teams who deserve better monitoring tools.*