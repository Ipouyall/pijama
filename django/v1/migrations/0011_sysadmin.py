# Generated by Django 4.2.3 on 2023-07-05 20:15

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    dependencies = [
        ('v1', '0010_extendeduser_chat_id'),
    ]

    operations = [
        migrations.CreateModel(
            name='SysAdmin',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('related_user', models.OneToOneField(on_delete=django.db.models.deletion.CASCADE, related_name='system_admin', to='v1.extendeduser', verbose_name='related_user')),
            ],
        ),
    ]
