#!/bin/sh
cd ddl;go generate;cd ..
cd entity;go generate;cd ..
cd enum;go generate;cd ..
cd repository;go generate;cd ..
cd vo;go generate;cd ..