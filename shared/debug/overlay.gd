extends Control


@onready var fps: Label = %FPSNumber
@onready var loglist: ItemList = %LogList
@onready var consoleedit: LineEdit = %ConsoleEdit


func _ready() -> void:
	pass # Replace with function body.


func _process(delta: float) -> void:
	fps.text = str(Performance.get_monitor(Performance.TIME_FPS))


func _unhandled_key_input(event: InputEvent) -> void:
	if event.is_action_released("overlay_toggle"):
		if self.visible:
			self.visible = false
			Input.mouse_mode = Input.MOUSE_MODE_CAPTURED
			consoleedit.grab_focus()
		else:
			self.visible = true
			Input.mouse_mode = Input.MOUSE_MODE_VISIBLE
			self.log("overlay revealed")


func log(text: String) -> void:
	loglist.add_item(text)


var expression = Expression.new()
func _on_console_edit_text_submitted(command: String) -> void:
	var error = expression.parse(command)
	if error != OK:
		loglist.add_item(expression.get_error_text())
		return

	consoleedit.clear()
	var result = expression.execute()
	if not expression.has_execute_failed():
		loglist.add_item(str(result))

