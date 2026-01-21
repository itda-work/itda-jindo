---
number: 27
title: Global/Local 경로 분리 지원
state: wip
labels:
  - enhancement
assignees: []
created_at: '2026-01-21T03:59:25Z'
updated_at: '2026-01-21T03:59:32Z'
---

## Background

현재 jindo CLI는 ~/.claude/ (Global) 경로만 지원합니다. Claude Code는 프로젝트별 .claude/ (Local) 경로도 공식 지원하므로, jindo도 두 경로를 모두 지원해야 합니다.

## Goal

1. **List 명령**: Global과 Local을 섹션으로 분리하여 출력
2. **CUD 명령**: --global / --local 플래그로 대상 경로 선택 (기본값: global)

## Non-goals

- monorepo 중첩 .claude/ 탐색
- git root 자동 탐색 (CWD 기준만 사용)

## Constraints

- 기존 동작과 backward compatible (플래그 없으면 Global)
- Local 경로가 없거나 비어있으면 해당 섹션 생략
