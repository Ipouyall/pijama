# Generated by Django 4.2.2 on 2023-07-06 15:39

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    dependencies = [
        ('v1', '0019_alter_extendeduser_city'),
    ]

    operations = [
        migrations.AlterField(
            model_name='package',
            name='city',
            field=models.ForeignKey(blank=True, null=True, on_delete=django.db.models.deletion.DO_NOTHING, to='v1.city'),
        ),
    ]
