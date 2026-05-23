import logging
from .SuperDAW import SuperDAW

def create_instance(c_instance):
    return SuperDAW(c_instance)
