package domain

// Relation which can exist between Class (and optionally further specified into Property)
type Relation struct {
	// Type of relation (free format)
	Type string

	// FromProperty can optionally be set if the Relation originates from a Property
	FromProperty *Property

	// From part of the Relation
	From *Class

	// ToProperty can optionally be set if the Relation points towards a Property (e.g. with an $anchor #)
	ToProperty *Property

	// To receiving end of the Relation
	To *Class
}
