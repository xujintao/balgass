#!/bin/bash
docker exec mailserver setup email add admin@r2f2.com 1234
docker exec mailserver setup email add noreply@r2f2.com 1234
docker exec mailserver setup alias add postmaster@r2f2.com admin@r2f2.com
docker exec mailserver setup alias add abuse@r2f2.com admin@r2f2.com
