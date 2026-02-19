# Database Migration Scripts for Windows
# Requires migrate.exe in PATH (from go install)

$DB_URL = "postgres://devops:devops123@localhost:5432/devops_db?sslmode=disable"
$MIGRATIONS_PATH = "./migrations"
$MIGRATE_CMD = "$env:USERPROFILE\go\bin\migrate.exe"

function Show-Help {
    Write-Host "Available commands:" -ForegroundColor Cyan
    Write-Host "  .\migrate.ps1 up          - Apply all pending migrations"
    Write-Host "  .\migrate.ps1 down        - Rollback the last migration"
    Write-Host "  .\migrate.ps1 create <name> - Create a new migration file"
    Write-Host "  .\migrate.ps1 version     - Show current migration version"
    Write-Host "  .\migrate.ps1 force <n>   - Force set migration version to n"
    Write-Host "  .\migrate.ps1 drop        - Drop everything (use with caution!)"
}

function Migrate-Up {
    Write-Host "Applying migrations..." -ForegroundColor Yellow
    & $MIGRATE_CMD -path $MIGRATIONS_PATH -database $DB_URL up
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Migrations applied successfully!" -ForegroundColor Green
    } else {
        Write-Host "✗ Migration failed!" -ForegroundColor Red
    }
}

function Migrate-Down {
    Write-Host "Rolling back last migration..." -ForegroundColor Yellow
    & $MIGRATE_CMD -path $MIGRATIONS_PATH -database $DB_URL down 1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Rollback successful!" -ForegroundColor Green
    } else {
        Write-Host "✗ Rollback failed!" -ForegroundColor Red
    }
}

function Migrate-Create {
    param([string]$Name)
    
    if ([string]::IsNullOrEmpty($Name)) {
        Write-Host "Error: Migration name is required!" -ForegroundColor Red
        Write-Host "Usage: .\migrate.ps1 create <migration_name>" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "Creating migration: $Name" -ForegroundColor Yellow
    & $MIGRATE_CMD create -ext sql -dir $MIGRATIONS_PATH -seq $Name
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Migration files created!" -ForegroundColor Green
    }
}

function Migrate-Version {
    Write-Host "Current migration version:" -ForegroundColor Yellow
    & $MIGRATE_CMD -path $MIGRATIONS_PATH -database $DB_URL version
}

function Migrate-Force {
    param([int]$Version)
    
    if ($Version -eq 0 -and $args[0] -ne "0") {
        Write-Host "Error: Version number is required!" -ForegroundColor Red
        Write-Host "Usage: .\migrate.ps1 force <version_number>" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "Forcing migration version to $Version..." -ForegroundColor Yellow
    & $MIGRATE_CMD -path $MIGRATIONS_PATH -database $DB_URL force $Version
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Version forced successfully!" -ForegroundColor Green
    }
}

function Migrate-Drop {
    $confirmation = Read-Host "WARNING: This will drop all tables! Are you sure? (yes/no)"
    if ($confirmation -eq "yes") {
        Write-Host "Dropping all tables..." -ForegroundColor Red
        & $MIGRATE_CMD -path $MIGRATIONS_PATH -database $DB_URL drop -f
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ All tables dropped!" -ForegroundColor Green
        }
    } else {
        Write-Host "Operation cancelled." -ForegroundColor Yellow
    }
}

# Main script logic
$command = $args[0]

switch ($command) {
    "up" { Migrate-Up }
    "down" { Migrate-Down }
    "create" { Migrate-Create -Name $args[1] }
    "version" { Migrate-Version }
    "force" { Migrate-Force -Version $args[1] }
    "drop" { Migrate-Drop }
    "help" { Show-Help }
    default { 
        if ([string]::IsNullOrEmpty($command)) {
            Show-Help
        } else {
            Write-Host "Unknown command: $command" -ForegroundColor Red
            Show-Help
        }
    }
}
