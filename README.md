# outputs
for a particular order of the outputs (from left to right), send the output commands to sway using the outputs resolution and scale. 

### Dependencies
[sway](https://github.com/swaywm/sway)

### Install
go 1.12

`go get github.com/nomoth/outputs` 

### Usage
`outputs --help`

### Sway config example for 2 external screens with a laptop with QHD resolution
```
set $output_cfg ~/go/bin/outputs HDMI-A-1 eDP-1 HDMI-A-2

bindswitch lid:on output eDP-1 disable
bindswitch lid:off output eDP-1 enable
bindswitch lid:toggle exec $output_cfg

output eDP-1 scale 1.6

exec_always $output_cfg
```

### Benefits
- always use the native resolution of the external screens whatever it is
- order is preserved even if a screen is disable or unavailable

### Limitation
Support only screens disposed on one line