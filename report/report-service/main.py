from app import create_app
import logging
import os
def main():
    create_app()

if __name__ == '__main__':
    logging.basicConfig(
        format='%(asctime)s.%(msecs)s:%(name)s:%(thread)d:' + '%(levelname)s:%(process)d:%(message)s',
        level=logging.INFO
    )
    print(os.getcwd())
    main()