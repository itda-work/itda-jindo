---
number: 19
title: install.sh 파일명 형식 수정 및 Windows 설치 스크립트 추가
state: done
labels:
  - bug
  - enhancement
assignees: []
created_at: '2026-01-20T14:55:38Z'
updated_at: '2026-01-20T14:55:45Z'
closed_at: '2026-01-20T14:55:45Z'
---

## 문제

- install.sh가 기대하는 파일명: jd_0.1.0_darwin_arm64.tar.gz
- 실제 릴리즈 파일명: jd-macos-arm64

## 해결

1. install.sh - 파일명 형식을 실제 릴리즈와 일치하도록 수정
2. install.ps1 - Windows용 설치 스크립트 추가
3. README.md - 플랫폼별 설치 안내 추가
