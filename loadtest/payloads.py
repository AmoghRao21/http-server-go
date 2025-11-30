from utils import rand_text, rand_int

def random_item():
    return {
        "name": rand_text(),
        "value": rand_int(),
        "tags": [rand_text(4), rand_text(4)],
        "meta": {"likes": rand_int(), "active": True}
    }

def partial_update():
    return {
        "value": rand_int()
    }
