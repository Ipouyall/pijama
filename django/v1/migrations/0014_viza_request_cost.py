# Generated by Django 4.2.2 on 2023-07-06 08:09

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('v1', '0013_alter_viza_serial_no'),
    ]

    operations = [
        migrations.AddField(
            model_name='viza',
            name='request_cost',
            field=models.DecimalField(decimal_places=0, default=500, max_digits=12),
            preserve_default=False,
        ),
    ]
