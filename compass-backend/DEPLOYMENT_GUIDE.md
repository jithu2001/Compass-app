# AWS EC2 Deployment Guide for Compass Backend

## Prerequisites
- AWS Account
- SSH key pair for EC2 access
- Domain name (optional, for production)

## Step 1: Launch EC2 Instance

### 1.1 Create EC2 Instance
1. Log in to AWS Console
2. Navigate to **EC2 Dashboard**
3. Click **Launch Instance**
4. Configure instance:
   - **Name**: `compass-backend-server`
   - **AMI**: Ubuntu Server 22.04 LTS (Free tier eligible)
   - **Instance Type**: t2.micro (free tier) or t2.small (recommended)
   - **Key pair**: Create new or select existing
   - **Network settings**:
     - Allow SSH (port 22) from your IP
     - Allow HTTP (port 80) from anywhere
     - Allow HTTPS (port 443) from anywhere
     - Allow Custom TCP (port 8080) from anywhere (for testing)
5. **Storage**: 20 GB gp3
6. Click **Launch Instance**

### 1.2 Configure Security Group
1. Go to **Security Groups** in EC2 Dashboard
2. Find your instance's security group
3. Add inbound rules:
   ```
   Type            Protocol    Port Range    Source
   SSH             TCP         22            Your IP/0.0.0.0/0
   HTTP            TCP         80            0.0.0.0/0
   HTTPS           TCP         443           0.0.0.0/0
   Custom TCP      TCP         8080          0.0.0.0/0
   PostgreSQL      TCP         5432          Same security group
   ```

## Step 2: Connect to EC2 Instance

```bash
# Get your instance public IP from AWS Console
ssh -i /path/to/your-key.pem ubuntu@<EC2-PUBLIC-IP>
```

## Step 3: Install Dependencies on EC2

### 3.1 Update System
```bash
sudo apt update && sudo apt upgrade -y
```

### 3.2 Install Go
```bash
# Download Go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz

# Extract and install
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

### 3.3 Install PostgreSQL
```bash
# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Verify installation
sudo systemctl status postgresql
```

## Step 4: Configure PostgreSQL

### 4.1 Create Database and User
```bash
# Switch to postgres user
sudo -u postgres psql

# In PostgreSQL prompt, run:
CREATE DATABASE compass_db;
CREATE USER compass_user WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE compass_db TO compass_user;
ALTER DATABASE compass_db OWNER TO compass_user;
\c compass_db
GRANT ALL ON SCHEMA public TO compass_user;
ALTER SCHEMA public OWNER TO compass_user;
\q
```

### 4.2 Configure PostgreSQL for Remote Access (if needed)
```bash
# Edit postgresql.conf
sudo nano /etc/postgresql/14/main/postgresql.conf
# Find and uncomment: listen_addresses = 'localhost'

# Edit pg_hba.conf
sudo nano /etc/postgresql/16/main/pg_hba.conf
# Add line: host    compass_db    compass_user    127.0.0.1/32    md5

# Restart PostgreSQL
sudo systemctl restart postgresql
```

## Step 5: Install Git
```bash
sudo apt install git -y
git --version
```

## Step 6: Clone Your Repository

### 6.1 Using HTTPS (recommended for private repos with token)
```bash
cd ~
git clone https://github.com/yourusername/Compass-backend.git
```

### 6.2 Configure Git (if private repo)
```bash
# Use GitHub Personal Access Token
git clone https://<YOUR_TOKEN>@github.com/yourusername/Compass-backend.git
```

## Step 7: Configure Application

### 7.1 Navigate to Project Directory
```bash
  cd ~/Compass-backend/compass-backend
```

### 7.2 Create Production Config
```bash
# Copy example config
cp config/config.example.yaml config/config.yaml

# Edit configuration
nano config/config.yaml
```

### 7.3 Update config.yaml for Production
```yaml
server:
  port: "8080"
  mode: "release"  # Change from debug to release

database:
  host: "localhost"
  port: "5432"
  user: "compass_user"
  password: "your_secure_password"
  name: "compass_db"

jwt:
  secret: "your-super-secret-jwt-key-change-this-in-production"
  access_token_duration: "15m"
  refresh_token_duration: "7d"
```

### 7.4 Install Go Dependencies
```bash
go mod download
go mod tidy
```

## Step 8: Build the Application

```bash
# Build binary
go build -o compass-backend cmd/api/main.go

# Make it executable
chmod +x compass-backend

# Test run (press Ctrl+C to stop)
./compass-backend
```

## Step 9: Setup Systemd Service (Production)

### 9.1 Create Service File
```bash
sudo nano /etc/systemd/system/compass-backend.service
```

### 9.2 Add Service Configuration
```ini
[Unit]
Description=Compass Backend API Service
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/Compass-app/compass-backend
ExecStart=/home/ubuntu/Compass-app/compass-backend/compass-backend
Restart=on-failure
RestartSec=10
StandardOutput=append:/var/log/compass-backend/output.log
StandardError=append:/var/log/compass-backend/error.log

Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
```

### 9.3 Create Log Directory
```bash
sudo mkdir -p /var/log/compass-backend
sudo chown ubuntu:ubuntu /var/log/compass-backend
```

### 9.4 Enable and Start Service
```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable compass-backend

# Start service
sudo systemctl start compass-backend

# Check status
sudo systemctl status compass-backend

# View logs
sudo journalctl -u compass-backend -f
```

## Step 10: Setup Nginx Reverse Proxy (Optional but Recommended)

### 10.1 Install Nginx
```bash
sudo apt install nginx -y
```

### 10.2 Configure Nginx
```bash
sudo nano /etc/nginx/sites-available/compass-backend
```

### 10.3 Add Nginx Configuration
```nginx
server {
    listen 80;
    server_name your-domain.com;  # Replace with your domain or EC2 IP

    client_max_body_size 10M;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 10.4 Enable Nginx Configuration
```bash
# Create symbolic link
sudo ln -s /etc/nginx/sites-available/compass-backend /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx
```

## Step 11: Setup SSL with Let's Encrypt (Optional)

### 11.1 Install Certbot
```bash
sudo apt install certbot python3-certbot-nginx -y
```

### 11.2 Obtain SSL Certificate
```bash
# Replace with your domain
sudo certbot --nginx -d your-domain.com

# Follow prompts to configure HTTPS
```

### 11.3 Auto-renewal Setup
```bash
# Test renewal
sudo certbot renew --dry-run

# Certbot automatically sets up a cron job for renewal
```

## Step 12: Initialize Database with Seed Data

### 12.1 Run Seed Script
```bash
cd ~/Compass-backend/compass-backend

# Make sure migrations have run (automatically happens on app start)
# Then run seed if you have one
go run db/seed.go
```

## Step 13: Firewall Configuration (UFW)

```bash
# Enable UFW
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 8080/tcp  # Backend (if accessing directly)
sudo ufw enable

# Check status
sudo ufw status
```

## Step 14: Testing Deployment

### 14.1 Test API
```bash
# From your local machine
curl http://<EC2-PUBLIC-IP>/auth/signin

# Or if using Nginx on port 80
curl http://<EC2-PUBLIC-IP>/auth/signin
```

### 14.2 Test from Frontend
Update your Flutter app's base URL to:
```dart
final baseUrl = 'http://<EC2-PUBLIC-IP>:80/api';
// or with domain
final baseUrl = 'https://your-domain.com/api';
```

## Step 15: Monitoring and Maintenance

### 15.1 View Application Logs
```bash
# Service logs
sudo journalctl -u compass-backend -f

# Application logs
tail -f /var/log/compass-backend/output.log
tail -f /var/log/compass-backend/error.log
```

### 15.2 Restart Service
```bash
sudo systemctl restart compass-backend
```

### 15.3 Update Application
```bash
# Pull latest code
cd ~/Compass-backend/compass-backend
git pull origin main

# Rebuild
go build -o compass-backend cmd/api/main.go

# Restart service
sudo systemctl restart compass-backend
```

## Step 16: Backup Strategy

### 16.1 Database Backup Script
```bash
# Create backup directory
mkdir -p ~/backups

# Create backup script
nano ~/backup-db.sh
```

Add to `backup-db.sh`:
```bash
#!/bin/bash
BACKUP_DIR="/home/ubuntu/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
pg_dump -U compass_user -h localhost compass_db > "$BACKUP_DIR/compass_db_$TIMESTAMP.sql"
# Keep only last 7 days of backups
find $BACKUP_DIR -name "compass_db_*.sql" -mtime +7 -delete
```

```bash
# Make executable
chmod +x ~/backup-db.sh

# Setup daily cron job
crontab -e
# Add line: 0 2 * * * /home/ubuntu/backup-db.sh
```

## Troubleshooting

### Service won't start
```bash
# Check logs
sudo journalctl -u compass-backend -n 50 --no-pager

# Check if port is in use
sudo lsof -i :8080

# Check permissions
ls -la /home/ubuntu/Compass-app/compass-backend/compass-backend
```

### Database connection issues
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Test connection
psql -U compass_user -h localhost -d compass_db

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-14-main.log
```

### Can't access from internet
```bash
# Check security group in AWS Console
# Check if service is running
sudo systemctl status compass-backend

# Check if port is open
sudo netstat -tulpn | grep 8080
```

## Security Best Practices

1. **Change default passwords** in config.yaml
2. **Use environment variables** for sensitive data
3. **Enable firewall (UFW)**
4. **Keep system updated**: `sudo apt update && sudo apt upgrade`
5. **Use SSL/TLS** with Let's Encrypt
6. **Regular backups** of database
7. **Monitor logs** regularly
8. **Use strong JWT secret**
9. **Limit SSH access** to specific IPs
10. **Set up CloudWatch** for monitoring (optional)

## Cost Optimization

- Use t2.micro for development/testing (free tier)
- Use t2.small or t3.small for production
- Consider RDS for database in production
- Use CloudFront CDN for static assets
- Set up auto-scaling for high traffic

## Next Steps

1. Set up monitoring with CloudWatch
2. Configure automated backups to S3
3. Implement CI/CD with GitHub Actions
4. Set up staging environment
5. Configure database replication for high availability

---

## Quick Reference Commands

```bash
# Service management
sudo systemctl start compass-backend
sudo systemctl stop compass-backend
sudo systemctl restart compass-backend
sudo systemctl status compass-backend

# View logs
sudo journalctl -u compass-backend -f

# Update application
cd ~/Compass-backend/compass-backend
git pull
go build -o compass-backend cmd/api/main.go
sudo systemctl restart compass-backend

# Database backup
pg_dump -U compass_user compass_db > backup.sql

# Database restore
psql -U compass_user compass_db < backup.sql
```
