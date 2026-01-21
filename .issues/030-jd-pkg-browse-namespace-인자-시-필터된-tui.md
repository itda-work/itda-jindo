---
number: 30
title: 'jd pkg browse: namespace 인자 시 필터된 TUI + Tab 자동완성'
state: done
labels:
  - enhancement
assignees: []
created_at: '2026-01-21T04:57:22Z'
updated_at: '2026-01-21T05:01:04Z'
closed_at: '2026-01-21T05:01:04Z'
---

## 문제

1. **TUI 문제**: `jd pkg browse`는 TUI가 작동하지만, `jd pkg browse verc-agen`처럼 namespace 인자를 주면 TUI 없이 CLI 출력만 됨

2. **자동완성 문제**: Tab 자동완성이 현재 파일/디렉토리로 되는데, 등록된 저장소 이름으로 자동완성 되어야 함

## 기대 동작

1. `jd pkg browse namespace` 실행 시 해당 namespace로 필터링된 TUI 표시
2. `jd pkg browse <Tab>` 시 등록된 저장소 이름으로 자동완성
