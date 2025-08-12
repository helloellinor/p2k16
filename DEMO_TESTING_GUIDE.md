# P2K16 Demo Mode Testing Guide

## Overview
This guide explains how to test the P2K16 system in demo mode and understand the enhanced logging features.

## Running Demo Mode

### Standalone Demo Server
```bash
go run cmd/demo/main.go
```

**Features:**
- ✅ Clean Bootstrap UI design
- ✅ Complete profile management system
- ✅ Password change functionality  
- ✅ Badge system
- ✅ Comprehensive logging
- ⚠️ No database required - all operations simulated

### Main Server (with Demo Fallback)
```bash
go run cmd/server/main.go
```

**Behavior:**
- Attempts database connection first
- Falls back to demo mode if database unavailable
- Uses custom design system (different from demo UI)
- All Phase 2 features available in both modes

## Enhanced Logging Features

### Demo Mode Logging
Look for these log indicators:
```
🎯 [DEMO MODE] timestamp | action | details
```

Examples:
- `🎯 [DEMO MODE] 20:01:22 | SERVER STARTUP | Demo server starting on port 8080`
- `🎯 [DEMO MODE] 20:01:43 | PAGE REQUEST | Home page visited`
- `🎯 [DEMO MODE] 20:01:43 | USER STATUS | Anonymous user`
- `🎯 [DEMO MODE] 20:02:15 | LOGIN SUCCESS | Demo user authenticated: demo`
- `🎯 [DEMO MODE] 20:02:30 | PASSWORD CHANGE SUCCESS | Password changed for user: demo`

### Handler Logging
For API operations:
```
🔧 [HANDLER] timestamp | action | details
```

Examples:
- `🔧 [HANDLER] 20:01:45 | PASSWORD CHANGE | Password change request received`
- `🔧 [HANDLER] 20:01:45 | VALIDATION ERROR | New passwords do not match`
- `🔧 [HANDLER] 20:01:50 | DEMO MODE | Password change simulated - no database update`

## Testing the Profile Management Features

### 1. Login to Demo Mode
- Navigate to `http://localhost:8080`
- Click "Login" 
- Use username: `demo`, password: `any`
- Watch the terminal for login logging

### 2. Access Profile Page
- Click "Profile" from dashboard
- Or navigate directly to `http://localhost:8080/profile`
- Observe the enhanced profile management UI

### 3. Test Password Change
- Fill out the password change form
- Try various validation scenarios:
  - Missing fields
  - Mismatched passwords
  - Short passwords
- Watch detailed logging for each validation step

### 4. Test Profile Updates
- Update name and phone number
- Submit the form
- Observe logging showing field updates

## Startup Messages

### Demo Mode
```
============================================================
🎭  P2K16 DEMO MODE - Development Testing Server
============================================================
📍 Server URL: http://localhost:8080
🔑 Demo Login: username='demo', password=any
📋 Features Available:
   • User authentication
   • Dashboard with badges
   • Profile management (password change, profile update)
   • Member listing
⚠️  Note: No database - all changes are simulated
============================================================
```

### Main Server (Demo Fallback)
```
============================================================
⚠️  P2K16 SERVER - FALLBACK TO DEMO MODE
============================================================
❌ Database connection failed: [error details]
🎭 Falling back to DEMO MODE - no database required
🔑 Demo logins available:
   • demo/password
   • super/super
   • foo/foo
⚠️  Note: All data operations will be simulated
============================================================
```

### Main Server (Production Mode)
```
============================================================
🚀 P2K16 SERVER - PRODUCTION MODE
============================================================
✅ Database connection successful
🗄️  Connected to: p2k16-web@localhost:2016/p2k16
💾 All data operations will be persisted to database
============================================================
```

## Key Improvements

1. **Clear Mode Indication**: Always know if you're in demo or production mode
2. **Comprehensive Logging**: Every user action is logged with context
3. **Missing Features Added**: Demo mode now has all Phase 2 profile management features
4. **Enhanced UX**: Clear visual feedback for all operations
5. **Better Testing**: Easy to see what's happening server-side during UI interactions

## Differences Between Demo and Main Server

| Feature | Demo Server | Main Server |
|---------|-------------|-------------|
| UI Design | Bootstrap CSS (CDN) | Custom design system |
| Database | None (simulated) | PostgreSQL or fallback |
| Routes | Simplified inline handlers | Full handler architecture |
| Logging | Detailed demo logging | Production + demo logging |
| Features | Phase 2 complete | Full feature set |

## Troubleshooting

### HTMX Not Working
If HTMX forms don't work due to CDN blocking:
- The forms will still work with standard HTTP POST
- Detailed logging will still show all server-side operations
- This is expected in sandboxed environments

### Database Connection Issues
The main server will automatically fall back to demo mode if database connection fails, providing a consistent testing experience.