extends CharacterBody3D


const SPEED = 5.0
const SPRINT_SPEED = 8.0
const JUMP_VELOCITY = 4.5
const MOUSE_SENSITIVITY = 0.001

# Get the gravity from the project settings to be synced with RigidBody nodes.
var gravity: float = ProjectSettings.get_setting("physics/3d/default_gravity")
var sprinting: bool = false
var crouching: bool = false
var speed: float = SPEED:
	set(val):
		pass
	get:
		if sprinting and not crouching:
			return SPRINT_SPEED
		else:
			return SPEED

@onready var camera: Camera3D = $Camera3D
@onready var collision: CollisionShape3D = $CollisionShape3D

func _ready() -> void:
	Input.mouse_mode = Input.MOUSE_MODE_CAPTURED


func _physics_process(delta: float) -> void:
	# Add the gravity.
	if not is_on_floor():
		velocity.y -= gravity * delta

	# Movement state changes are only valid while on the floor
	# Sorry no crouch-jumps or funky air sprinting
	if is_on_floor():
		# Sprinting
		if Input.is_action_just_pressed("player_sprint"):
			sprinting = !sprinting

		# Crouching
		if Input.is_action_just_pressed("player_crouch"):
			if crouching:
				crouching = false
				collision.shape.height = 2
			else:
				crouching = true
				# Shrinking the collision shape changes the size of the
				# overall player which will move the camera too
				collision.shape.height = 1

		# Handle jump.
		if Input.is_action_just_pressed("player_jump"):
			velocity.y = JUMP_VELOCITY

	# Get the input direction and handle the movement/deceleration.
	# As good practice, you should replace UI actions with custom gameplay actions.
	var input_dir := Input.get_vector("player_left", "player_right", "player_forward", "player_backward")
	var direction := (transform.basis * Vector3(input_dir.x, 0, input_dir.y)).normalized()
	if direction:
		velocity.x = direction.x * SPEED
		velocity.z = direction.z * SPEED
	else:
		velocity.x = move_toward(velocity.x, 0, SPEED)
		velocity.z = move_toward(velocity.z, 0, SPEED)

	move_and_slide()


func _unhandled_key_input(event: InputEvent) -> void:
	if event.is_action_released("mouse_toggle"):
		if Input.mouse_mode == Input.MOUSE_MODE_CAPTURED:
			Input.mouse_mode = Input.MOUSE_MODE_VISIBLE
		else:
			Input.mouse_mode = Input.MOUSE_MODE_CAPTURED


func _unhandled_input(event: InputEvent) -> void:
	if event is InputEventMouseMotion and Input.mouse_mode == Input.MOUSE_MODE_CAPTURED:
		var look_dir = event.relative * MOUSE_SENSITIVITY
		# Horizontal rotation moves the whole player
		self.rotation.y -= look_dir.x
		# Vertical rotation only moves the camera and is clamped
		camera.rotation.x = clampf(camera.rotation.x - look_dir.y, -1.5, 1.5)

