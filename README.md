# OHLC Price Data

## Table of Contents
* [General Info](#general-info)
* [Installation & Build](#installation--build)

## General Info
We have just purchased a large amount of historical OHLC price data that were shared to us in CSV files
format. We need to start centralising and digitalising those data. These files can be ranging from a few GBs to
a couple of TBs.

## Installation & Build
1. First of All, try to install docker. you can install docker by this [guideline](https://docs.docker.com/engine/install/)
2. Make Sure 8080 and 5432 ports are open, Then run application and Database container by `docker-compose up -d --build`
3. Open your browser, after running previous step you should see swagger: `localhost:8080/swagger/index.html`
