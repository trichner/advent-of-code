import functools
import json


def parse_input(fname):
  with open(fname, 'r') as file:
    lines = [l.strip() for l in file.readlines()]
    lines = filter(lambda l: len(l) > 0, lines)
    return list(map(lambda l: json.loads(l), lines))


def group(lines) -> list:
  groups = []
  group = []
  for i, el in enumerate(lines):
    if i % 2 == 0:
      group = [el]
      groups.append(group)
    else:
      group.append(el)

  return groups


def cmp(l, r):
  if type(l) == int and type(r) == int:
    return l - r

  if type(l) == int and type(r) == list:
    return cmp([l], r)

  if type(l) == list and type(r) == int:
    return cmp(l, [r])

  for i in range(max(len(l), len(r))):
    if i >= len(l):
      return -1
    if i >= len(r):
      return 1

    c = cmp(l[i], r[i])
    if c != 0:
      return c

  return 0


def part_one():
  groups = group(parse_input('input.txt'))
  total = 0
  for i, g in enumerate(groups):
    l = g[0]
    r = g[1]
    e = cmp(l, r)
    if e == 0:
      raise 'oops'
    if e < 0:
      total += i + 1

  print(total)


def part_two():
  elements = parse_input('input.txt')
  elements.append([[2]])
  elements.append([[6]])

  elements = sorted(elements, key=functools.cmp_to_key(cmp))

  total = 1
  for i, e in enumerate(elements):
    if type(e) == list and len(e) == 1 and type(e[0]) == list and len(e[0]) == 1:
      v = e[0][0]
      if v == 6 or v == 2:
        total *= i + 1

  print(total)


if __name__ == '__main__':
  part_one()
  part_two()
