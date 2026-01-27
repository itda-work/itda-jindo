---
number: 47
title: Windows에서 os.UserConfigDir() 사용하여 config 경로 결정
state: done
labels:
  - enhancement
assignees: []
created_at: '2026-01-27T05:58:11Z'
updated_at: '2026-01-27T06:03:34Z'
closed_at: '2026-01-27T06:03:34Z'
---

## 개요

현재 `pkg/config/paths.go`의 `GetConfigDir()`는 모든 OS에서
`$HOME/.config/itda-skills/` 경로를 사용한다.
Windows에서는 관례에 맞지 않으므로,
Windows만 `os.UserConfigDir()`를 사용하도록 변경한다.

## 현재 동작

모든 OS에서 동일한 로직:

- `$XDG_CONFIG_HOME/itda-skills/` (환경변수 설정 시)
- `$HOME/.config/itda-skills/` (기본값)

## 변경 후 동작

| OS      | 경로                                                   |
| ------- | ------------------------------------------------------ |
| Linux   | `~/.config/itda-skills/` (변경 없음)                   |
| macOS   | `~/.config/itda-skills/` (변경 없음)                   |
| Windows | `%AppData%\\itda-skills\\` (`os.UserConfigDir()` 사용) |

## 구현 방향

- `pkg/config/paths.go`의 `GetConfigDir()` 함수 수정
- `runtime.GOOS == "windows"`일 때 `os.UserConfigDir()` 사용
- Linux/macOS는 기존 XDG 스타일 로직 유지
- Windows 사용자가 현재 없으므로 마이그레이션 로직 불필요

## 의존 프로젝트 영향 분석

아래 3개 프로젝트가 `github.com/itda-skills/jindo/pkg/config`를 직접 import하여 사용 중이다.

### 1. skills.law

- **import**: `jdconfig "github.com/itda-skills/jindo/pkg/config"`
- **사용 함수**: `Load()`, `GetConfigDir()`, `GetConfigPath()`, `InitConfig()`
- **파일**: `internal/config/config.go`
- **하드코딩 경로**: `~/.config/itda-skills/law/metrics` (line 129)

### 2. skills.quant-data

- **import**: `jinconfig "github.com/itda-skills/jindo/pkg/config"`
- **사용 함수**: `Load()`, `GetConfigPath()`, `ConfigExists()`, `InitConfig()`
- **파일**: `internal/config/config.go`, `internal/cli/config.go`
- **하드코딩 경로**: CLI 도움말에 `~/.config/itda-skills/config.toml` (line 26)

### 3. skills.tts

- **import**: `"github.com/itda-skills/jindo/pkg/config"`
- **사용 함수**: `Load()`, `GetWithEnv()`
- **파일**: `internal/config/config.go`
- **하드코딩 경로**: 주석에 `~/.config/itda-skills/config.toml` (line 3)

### 영향 범위

- **코드 변경 불필요**: 3개 프로젝트 모두 jindo의
  `GetConfigDir()`/`GetConfigPath()`를 호출하므로,
  jindo 수정 후 재빌드하면 자동으로 Windows 경로가 적용됨
- **하드코딩 경로 수정 필요**: CLI 도움말/주석에
  `~/.config/itda-skills/` 경로가 하드코딩된 부분은
  OS별 동적 표시 또는 문구 수정 검토 필요
  - `skills.law`: `internal/config/config.go:129`
  - `skills.quant-data`: `internal/cli/config.go:26`
  - `skills.tts`: `internal/config/config.go:3`
