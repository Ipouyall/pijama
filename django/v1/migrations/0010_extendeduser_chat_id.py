# Generated by Django 4.2.2 on 2023-07-05 18:26

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('v1', '0009_rename_reserved_hotel_id_treatmentrequest_reserved_hotel'),
    ]

    operations = [
        migrations.AddField(
            model_name='extendeduser',
            name='chat_id',
            field=models.IntegerField(default=0),
            preserve_default=False,
        ),
    ]
