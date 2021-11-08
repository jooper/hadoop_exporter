#!/usr/bin/env bash

nohup ./nm_exporter 2>&1 &
tail -f nohup.out