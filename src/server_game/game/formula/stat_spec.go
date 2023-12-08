package formula

func StatSpec_GetPercent(class, id, strengh, dexterity, vitality, energy, leadership int, percent *float32) {
	call(f.StatSpec, "StatSpec_GetPercent", "iiiiiii>d", class, id, strengh, dexterity, vitality, energy, leadership, percent)
}
