#!/bin/bash
# Setup CBP Scheduler for auto-start on boot/wake

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STARTUP_SCRIPT="$SCRIPT_DIR/start_scheduler.sh"

echo "==========================================="
echo "CBP SCHEDULER AUTO-START SETUP"
echo "==========================================="

# Option 1: Add to .bashrc for user login
echo ""
echo "Setting up .bashrc auto-start..."
BASHRC_MARKER="# CBP_SCHEDULER_AUTOSTART"

if grep -q "$BASHRC_MARKER" ~/.bashrc; then
    echo "  Already configured in .bashrc"
else
    cat >> ~/.bashrc << EOF

$BASHRC_MARKER
# Auto-start CBP Scheduler on login
if [ ! -f "\$HOME/.cbp_scheduler/scheduler.pid" ] || ! ps -p \$(cat "\$HOME/.cbp_scheduler/scheduler.pid" 2>/dev/null) > /dev/null 2>&1; then
    bash $STARTUP_SCRIPT > /dev/null 2>&1
fi
EOF
    echo "  ✅ Added to .bashrc"
fi

# Option 2: Create systemd user service (if systemd is available)
if command -v systemctl > /dev/null 2>&1; then
    echo ""
    echo "Setting up systemd user service..."

    SYSTEMD_USER_DIR="$HOME/.config/systemd/user"
    mkdir -p "$SYSTEMD_USER_DIR"

    cat > "$SYSTEMD_USER_DIR/cbp-scheduler.service" << EOF
[Unit]
Description=CBP Society Autonomous Scheduler
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/python3 $SCRIPT_DIR/cbp_scheduler_daemon.py
ExecStop=/usr/bin/python3 $SCRIPT_DIR/cbp_scheduler_daemon.py stop
Restart=always
RestartSec=30

[Install]
WantedBy=default.target
EOF

    # Reload and enable
    systemctl --user daemon-reload
    systemctl --user enable cbp-scheduler.service
    echo "  ✅ Systemd service created and enabled"
    echo "     Start now: systemctl --user start cbp-scheduler"
    echo "     Check status: systemctl --user status cbp-scheduler"
fi

# Option 3: Add to crontab for reboot
echo ""
echo "Setting up crontab @reboot..."
CRON_CMD="@reboot sleep 30 && bash $STARTUP_SCRIPT"

if crontab -l 2>/dev/null | grep -q "cbp-chain/start_scheduler"; then
    echo "  Already in crontab"
else
    (crontab -l 2>/dev/null; echo "$CRON_CMD") | crontab -
    echo "  ✅ Added to crontab"
fi

echo ""
echo "==========================================="
echo "AUTO-START SETUP COMPLETE!"
echo "==========================================="
echo ""
echo "The scheduler will now start automatically on:"
echo "  - User login (.bashrc)"
echo "  - System reboot (crontab)"
if command -v systemctl > /dev/null 2>&1; then
    echo "  - As systemd service"
fi
echo ""
echo "Manual controls:"
echo "  Start: bash $SCRIPT_DIR/start_scheduler.sh"
echo "  Status: python3 $SCRIPT_DIR/cbp_scheduler_daemon.py status"
echo "  Stop: python3 $SCRIPT_DIR/cbp_scheduler_daemon.py stop"
echo ""
echo "Logs: ~/.cbp_scheduler/logs/"