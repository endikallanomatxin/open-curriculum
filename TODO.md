- Manifest:

    - Repasar traducción a euskera

    - Explicar que ocurre cuando se modifican unidades sobre las que ibas a hacer una propuesta.

        Cuando empiezas a hacer una propuesta, el curriculum debería congelarse. Luego si mientras que la preparas surgen modificaciones, puedes ir (resolviendo los merge conflicts) incorporando los cambios antes de presentar tu propuesta.
    
    - Redo graph images

    - Redo proposal and poll images

- Create an explainer video

- Platform demo:
    - Users
    - Teach:
        - Proposals
            - [x] Creation
            - [x] Activation
            - [x] Submission
            - [/] Adding changes:
                - [x] Unit creation
                - [x] Unit deletion
                - [x] Unit rename
                - [ ] Unit content modification
                - [x] Dependency creation
                - [x] Dependency deletion
                - [ ] Transfer certifications and reads (it has to be done every time you delete something)
            - Fix:
                - Cuando se intenta abrir una unidad que no existe (que puede pasar porque queda guardado en las cookies), se para el programa.
                - When you click the close button of the dialogue, the page should be refreshed.
                - Al crear una propuesta que se autoactive.
        - Voting:
            - [x] Basic single proposal poll voting
            - [ ] Multiple proposals poll voting
            - [ ] See relevant polls
            - [ ] Open a poll and vote on it
    - Learn
        - [x] See tree
        - [x] See unit details
        - [ ] See unit documents, tasks...
        - [ ] Mark units as seen.
        - [ ] Unit recommendation engine

- CSRF protection



Por decidir:

- Las unidades tienen descripción?

- Las unidades pueden tener más de un documento / video ? Si solo tienen una se podría simplificar las operaciones.

- Las dependencias son independientes (se crean y destruyen) o son una propidad de una unidad (y se modifica la unidad)?

- Can unstaged changes exist?

- Can a unit have multiple documents or explainer videos?

- Maybe dependencies need to have a description explaining why that dependency exist, to help future modifications


Future features:

- Comments for improvement requests.