---
number: 29
title: global/local 플래그 상호 배타적 검증 추가
state: done
labels:
  - enhancement
assignees: []
created_at: '2026-01-21T04:47:01Z'
updated_at: '2026-01-21T13:13:50Z'
closed_at: '2026-01-21T13:13:50Z'
---

skills, agents, commands, hooks의 CUD 명령에서 --global과 --local 플래그가 동시에 지정되면 에러를 반환해야 합니다.
