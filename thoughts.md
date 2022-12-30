# examples

## dayz-style

- insulation: static value (0 -> 1)
- external temperature affects body temperature: effect is one-way (body won't appreciably affect environment). temperature drop (rate) is proportional to the temperature difference, multiplied by (1-insulation)
- body temperature consumes food: the lower the body temp below target, consumption is square of this. for each unit of food consumed, increase temperature linearly
- health: increases proportionally to blood level. decreases proportionally to the square of the absolute body temperature difference from ideal.
- blood: consumes food and water. rate of increase is proportional to health (but consumption is the same)
- stamina: consumes food and water linearly. rate of increase is a ~normal distribution: fastest in the middle

# components

- conversion factor (how many units of one thing translates to how many of another). Often linear, sometimes might not be
- conversion rate (how quickly does the conversion occur?)

# considerations

- not a DAG; can have bidirectional effects.
- most effects are "consumption", meaning source value(s) are modified, not just the target

# approach

- "linker" is a function strongly associated with the "target", and has 0+ sources. [e.g. food+water are sources for blood]

- create a person
- assign attributes
- link attributes: each association will have its own custom linker implementation, but could drop on common helper functions. DRY later maybe, if it looks like helping

## example

(assuming all values are between 0 and 1)

### blood: (food, water, health)

conversion factor: 1 food + 1 water -> (health) blood
rate: 0.01 blood/minute [e.g. at 50% health, 0.01 food + 0.01 water -> 0.005 blood per minute)

### stamina: (food, water)

conversion factor: 1 food + 1 water -> 1 stamina + 0.1 heat
rate: (1 - abs(stamina-0.5)) \* 0.1 / minute
e.g. at 50% stamina, rate is 0.1/minute. at 10% (or 90%) stamina, rate is 0.06/minute)

### body temp: (env. temp, food, water)

#### if hot:

conversion factor: 1 water -> -(1-insulation) temperature
rate: (body temp-target)^2/minute

#### if cold:

conversion factor: 1 food -> 1 temperature
rate: (target - body temp)^2/minute
