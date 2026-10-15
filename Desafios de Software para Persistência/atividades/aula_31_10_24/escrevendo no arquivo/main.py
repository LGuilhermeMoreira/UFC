with open('file.txt','w',encoding='utf-8') as file:
    while True: 
        try:
            entrada = input('digite algo: ')
            if entrada == '':
                break
            print(entrada,file=file)
        except EOFError:
           print(EOFError)
           break