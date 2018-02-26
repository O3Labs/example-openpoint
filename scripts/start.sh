#!/bin/bash
cp /var/github.com/openpoint/systemd/openlotalty.service /etc/systemd/system/openlotalty.service
sudo systemctl daemon-reload
sudo systemctl enable openlotalty
sudo systemctl restart openlotalty
exit $?	