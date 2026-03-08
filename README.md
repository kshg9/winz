# winz

`winz` is a lightweight, single-binary Go CLI that scaffolds semester-lab exercise folders from embedded templates.

## Build

```bash
go build -o tool .
```

The binary embeds `internal/templates/**`, so no external template files are needed at runtime.

## CLI Usage

```bash
tool list
tool init <template>
tool init <template> <target-dir>
tool tui
```

Examples:

```bash
tool list
tool init dummy_exercise
tool init bi_lab/exp1_schema_design ./my-bi-exp
```

### Commands

- `list`: print all available embedded template paths.
- `init`: generate template files into the target directory (or the template basename by default).
- `tui`: open the minimal interactive picker with two modes: normal navigation (`j/k`) and search mode (`f` to enter, `esc` to exit). Large template sets are windowed so only a slice is rendered on-screen while you scroll.
- `uninstall`: run existing platform uninstall behavior.

## Embedded Template Layout

Add new templates under `internal/templates/`:

```text
internal/templates/
  <lab_name>/
    <experiment_name>/
      README.md
      ...any files/folders
```

The CLI auto-detects templates without additional code changes. It discovers top-level experiment roots that contain files, so subject folders can hold many experiments (including deeper paths like `web_lab/angular/exp_component`).

## Rendering Rules

- Files ending in `.tmpl` are rendered and written without the `.tmpl` suffix.
- All other files are copied as-is.
- On Windows, generated files are normalized to CRLF newlines.

## Included 6th Semester IT Lab Dummy Templates

- `bi_lab/exp1_schema_design`
- `web_lab/exp_typescript_basic`
- `web_lab/exp_advanced_ts`
- `web_lab/angular/exp_component_basics`
- `sensor_lab/exp_sensor_interface`
- `mad_pwa_lab/exp_flutter_ui`
- `ds_python_lab/exp_numpy_pandas`
- `rest_api_lab/exp_go_rest_api`
- `dummy_exercise`
