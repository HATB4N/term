import sys
import argparse
import subprocess

def main():
    print('대충 상태메세지')
    print('type -h or --help for help(wip)\n')
    parser = argparse.ArgumentParser()
    parser.add_argument('-r', '--range', dest='range', action='store')

    args = parser.parse_args()

    call = ['./scan', args.range]

    res = subprocess.run(call)
    print(res.stdout)


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print('인자 줘')
    else:
        main()