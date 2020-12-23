#!/usr/bin/env python
import os
from neo4j import debug
import neo4j

debug.watch("neo4j")

password = os.environ.get("NEO4J_PASSWORD", "password")

with neo4j.GraphDatabase.driver("bolt://localhost:8888", auth=("neo4j", password)) as driver:
    for i in range(1, 3):
        with driver.session(database="neo4j") as s:
            r = s.run('return 1')
            for i in r:
                print(i)
