package main

const crewRegex = `(?<=\()\w+(?:,\s*\w+(?:\s+\w+)*)*(?=\))`

const UnitRegex = `[A-Z]+(?:\/[A-Z]+)?-\d+[A-Z]?`

const craftRegex = `(?:AH-94|AV-42C|F-45A|F\/A-26B|EF-24G|T-55)`

const id64RegEx = `\d{17}`
