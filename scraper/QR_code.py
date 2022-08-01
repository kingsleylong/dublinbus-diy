import uuid
import qrcode

def QrCode(input_data):
        #Creating an instance of qrcode
        qr = qrcode.QRCode(
            version=1,
            box_size=10,
            border=5)
        qr.add_data(input_data)
        qr.make(fit=True)
        img = qr.make_image(fill='black', back_color='white')
        filename = img.save('qr_code_' + str(uuid.uuid4()) + '.png')
        return filename

# I put the link to our desktop app, and 2 links to App and Play store 
# Todo: change links when the app is up on the store
qr_code_list = ["http://ipa-003.ucd.ie/#/","https://apps.apple.com/ie/app/planny-day-planner/id1515324201", "https://play.google.com/store/apps/details?id=com.thomasyoung.dadish3"]

for line in qr_code_list:
    QrCode({line})

