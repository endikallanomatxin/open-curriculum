## Graph

### Properties

Unit have the following relevance metrics:

- Popularity. Measured in views.

- Pay. Measured in dollars.

- Reaching goals.

All of them are condensed into a relevance metric.


### Rendering

#### Order. Vertical position

Accumulated relevance is the sum of all the relevances of the units that depend on that unit (directly or indirectly).

To order the graph, use a Kahn algorithm, but each step, the unit with the highest accumulated relevance is selected.

#### Horizontal position

Hay que encontrar una forma de elegir la posición horizontal que minimice los cruces. Por ahora vamos a hacer:

- Buscar todas las leaves y asignarles una posición horizontal linspace de 0 a 1.
- Para cada nodo arriba en el grafo, asignarle la posición horizontal promedio de sus hijos.




## Groups

- Units cannot depend on groups, groups are just for contextualizing.

