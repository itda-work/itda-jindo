---
number: 25
title: jd list 명령 추가 및 모든 서브 명령에 축약 지원
state: done
labels:
  - enhancement
assignees: []
created_at: '2026-01-21T02:12:27Z'
updated_at: '2026-01-21T02:17:00Z'
closed_at: '2026-01-21T02:17:00Z'
---

## 개요

`jd list` 또는 `jd l` 명령으로 모든 skills, agents, commands, hooks 목록을 출력하는 기능 추가.

## 세부 요구사항

### 1. jd list / jd l 명령

- 모든 skills, agents, commands, hooks 목록 출력
- 카테고리별 그룹핑 형식

### 2. 모든 서브 명령에 축약 지원

- 범위: 모든 서브 명령 (list, add, remove, show 등)
- 충돌 처리: 선착순 우선 (먼저 등록된 명령이 짧은 축약 소유)
- 예시: `jd s list` → `jd s l`, `jd s add` → `jd s a`
