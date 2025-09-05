#!/bin/bash

# SAMS Database Backup Script
# This script creates a backup of the SAMS PostgreSQL database

# Configuration
BACKUP_DIR="/home/ubuntu/SAMS/backup"
DB_CONTAINER="sams-postgres"
DB_USER="sams_user"
DB_NAME="sams_db"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="sams_db_backup_${TIMESTAMP}.sql"
COMPRESSED_FILE="sams_db_backup_${TIMESTAMP}.sql.gz"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

echo "Starting SAMS database backup..."
echo "Timestamp: $TIMESTAMP"
echo "Backup file: $BACKUP_FILE"

# Create database backup
echo "Creating database dump..."
docker exec $DB_CONTAINER pg_dump -U $DB_USER -d $DB_NAME > "$BACKUP_DIR/$BACKUP_FILE"

# Check if backup was successful
if [ $? -eq 0 ]; then
    echo "✅ Database backup completed successfully!"
    echo "Backup file: $BACKUP_DIR/$BACKUP_FILE"
    
    # Create compressed version
    echo "Creating compressed backup..."
    gzip -c "$BACKUP_DIR/$BACKUP_FILE" > "$BACKUP_DIR/$COMPRESSED_FILE"
    
    if [ $? -eq 0 ]; then
        echo "✅ Compressed backup created: $BACKUP_DIR/$COMPRESSED_FILE"
    fi
    
    # Show backup file sizes
    echo ""
    echo "Backup file sizes:"
    ls -lah "$BACKUP_DIR/$BACKUP_FILE" "$BACKUP_DIR/$COMPRESSED_FILE" 2>/dev/null
    
    # Clean up old backups (keep last 7 days)
    echo ""
    echo "Cleaning up old backups (keeping last 7 days)..."
    find "$BACKUP_DIR" -name "sams_db_backup_*.sql" -mtime +7 -delete
    find "$BACKUP_DIR" -name "sams_db_backup_*.sql.gz" -mtime +7 -delete
    
    echo "✅ Backup process completed successfully!"
    
else
    echo "❌ Database backup failed!"
    exit 1
fi
