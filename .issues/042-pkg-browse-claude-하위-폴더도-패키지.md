---
number: 42
title: 'pkg browse: .claude/ 하위 폴더도 패키지 검색 대상에 추가'
state: done
labels:
  - enhancement
assignees: []
created_at: '2026-01-23T05:16:15Z'
updated_at: '2026-01-23T05:18:27Z'
closed_at: '2026-01-23T05:18:27Z'
---

## 개요

패키지 브라우징 시 루트 레벨 폴더뿐만 아니라 `.claude/` 하위 폴더도 검사하도록 개선

## 현재 상태

- 패키지 브라우징 시 루트 레벨 폴더만 검사
  - `skills/`
  - `commands/`
  - `agents/`
  - `hooks/`

## 변경 사항

- `.claude/` 하위 폴더도 추가 검사
  - `.claude/skills/`
  - `.claude/commands/`
  - `.claude/agents/`
  - `.claude/hooks/`
- 동일 이름의 패키지가 양쪽에 있으면 경로로 구분하여 둘 다 표시

## 영향 범위

- `internal/pkg/repo/repo.go`의 스캔 함수들:
  - `scanSkills`
  - `scanCommands`
  - `scanAgents`
  - `scanHooks`
